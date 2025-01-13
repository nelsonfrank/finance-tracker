import { Axios } from "@/utils/axiosClient";
import { AxiosResponse } from "axios";
import { User } from "./users";

export interface loginPayload {
  email: string;
  password: string;
}
export interface loginResponse {
  user: User;
  access_token: string;
  refresh_token: string;
}
export interface registerPayload {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
}
export type registerResponse = User;

export const loginAPI = (
  payload: loginPayload
): Promise<AxiosResponse<loginResponse>> =>
  Axios.post("v1/auth/login", payload);

export const registerUserAPI = (
  payload: registerPayload
): Promise<AxiosResponse<registerResponse>> =>
  Axios.post(`/v1/auth/register`, payload);
