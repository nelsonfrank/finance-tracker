import { createEnv } from "@t3-oss/env-nextjs";
import { z } from "zod";

export const env = createEnv({
  server: {
    NODE_ENV: z
      .enum(["development", "test", "production"])
      .default("development"),
    API_BASE_URL: z.string().url(),
    APP_BASE_URL: z.string().url(),
  },

  experimental__runtimeEnv: process.env,
});
