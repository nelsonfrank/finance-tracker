import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { fetchUsers, User } from "@/data/backend/api/users";

export const useUsers = (): UseQueryResult<User[], Error> => {
  return useQuery({
    queryKey: ["users"],
    queryFn: fetchUsers,
    staleTime: 1000 * 60 * 5,
    refetchOnWindowFocus: false,
  });
};
