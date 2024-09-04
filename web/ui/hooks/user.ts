import client from "@/lib/client";
import useAuthStore from "@/store/auth";
import { useQuery } from "@tanstack/react-query";

export const useInitializeUserData = () => {

  const { setUser } = useAuthStore()

  const { data, isError } = useQuery({
    queryKey: ['userData'],
    queryFn: () => client.user.userList(),
    staleTime: Infinity, // Prevent auto-refetching
    retry: false, // Disable retries
  });

  // Properly handle this case 
  if (data?.data === undefined || isError) {
    return
  }

  setUser(data.data.user)
};
