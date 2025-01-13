import { fetcher } from "@/utils/apiClient";

export interface User {
  ID: number;
  first_name: string;
  last_name: string;
  email: string;
}

export const fetchUsers = (): Promise<User[]> => fetcher("/api/users");
export const fetchUserById = (id: number): Promise<User> =>
  fetcher(`/api/users/${id}`);
