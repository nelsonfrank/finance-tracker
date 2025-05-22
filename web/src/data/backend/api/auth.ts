import { Axios } from "@/utils/axiosClient";
import { AxiosResponse } from "axios";
import { User } from "./users";
import { tryCatchWrapper } from "@/utils/api-utils";

export interface loginPayload {
  email: string;
  password: string;
}
export interface loginResponse {
  user: User;
  access_token: string;
  access_token_expires: string;
  refresh_token: string;
}
export interface registerPayload {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
}

export type registerResponse = User;

export interface refreshTokenDTO {
  access_token: string;
  access_token_expires: string;
  refresh_token: string
}

export const loginAPI = (payload: loginPayload): Promise<AxiosResponse<loginResponse>> => {
  return tryCatchWrapper<AxiosResponse<loginResponse>>(async () => {
    const data = await Axios.post("v1/auth/login", payload);
    return data
  });
};

export const registerUserAPI = (
  payload: registerPayload
): Promise<AxiosResponse<registerResponse>> =>
  Axios.post(`/v1/auth/register`, payload);

export const refreshTokenAPI = (payload: { refresh_token: string }): Promise<AxiosResponse<refreshTokenDTO>> =>
  Axios.post(
    "/v1/auth/refresh-token",
    payload
  );