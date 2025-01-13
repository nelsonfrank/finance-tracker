import { Axios } from "@/utils/axiosClient";
import { AxiosResponse } from "axios";

export interface dashbaordResponse {
  totalMonthlyIncome: string;
}
export const userDashboardAPI = (): Promise<AxiosResponse<dashbaordResponse>> =>
  Axios.get("v1/dashboard");
