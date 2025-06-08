import { loginAPI, refreshTokenAPI } from "@/data/backend/api"
import type { DefaultSession, DefaultUser, NextAuthOptions } from "next-auth"
import { JWT } from "next-auth/jwt"
import Credentials from "next-auth/providers/credentials"


declare module "next-auth" {
  interface Session extends DefaultSession {
    user: User & DefaultSession["user"]
  }

  interface User extends DefaultUser {
    id: string
    email: string
    name: string
    access_token: string
    access_token_expires: string
    refresh_token: string
  }

  
}

declare module "next-auth/jwt" {
  interface JWT {
    sub?: string
    access_token?: string
    access_token_expires: string
    refresh_token?: string
  }
}

async function refreshAccessToken(token: JWT) {
  try {
   const response = await refreshTokenAPI({refresh_token: token.refresh_token ?? ""})
    
    const { access_token, access_token_expires, refresh_token }= response.data;
    console.log({response})
    return {
      ...token,
      access_token: access_token,
      access_token_expires: access_token_expires,
      refresh_token: refresh_token ?? token.refresh_token, // Fall back to old refresh token
    }
  } catch (error) {
    console.log(error)

    return {
      ...token,
      error: "RefreshAccessTokenError",
    }
  }
}

export const authConfigs = {
  providers: [
    Credentials({
      name: "credentials",
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" }
      },
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) return null
          try {
            const {email, password} = credentials;
            const {data} = await loginAPI({email, password})
            
            if (!data) return null
          
            const {user, access_token, refresh_token, access_token_expires} = data;

            return {
                id: String(user.ID),
                email: user.email,
                name: `${user.first_name} ${user.last_name}`,
                access_token: access_token,
                access_token_expires: access_token_expires,
                refresh_token: refresh_token
            }
          } catch (error) {
              if (error instanceof Error) {
                throw new Error(error.message || "Something went wrong");
              }
              throw new Error("An unexpected error occurred");
          }
      }
    })
  ],
  pages: {
    signIn: "/auth/login",
    error: "/auth/login",
    signOut: "/",
  },
  callbacks: {
      jwt({ token, user }) {
        if (user) {
          token.sub = user.id
          token.access_token = user.access_token
          token.access_token_expires = user.access_token_expires ?? ''
          token.refresh_token = user.refresh_token
          

          return token
        }

      
        const tokenExpiresIn = new Date(token.access_token_expires);
        const now = new Date();
        
        if(now < tokenExpiresIn){
          return token
        }

        return refreshAccessToken(token)
      },
      async session({ session, token }) {
        if (session.user) {
          session.user.id = token.sub as string
          session.user.access_token = token.access_token as string
          session.user.refresh_token = token.refresh_token as string;
          session.user.access_token_expires = token.access_token_expires;
        }
        return session
      }
  }
}satisfies NextAuthOptions

