import axios from "axios";



export const tryCatchWrapper = async <T>(fn: () => Promise<T>): Promise<T> => {
    try {
      return await fn();
    } catch (error) {
        if (axios.isAxiosError(error)) {
            if (error.response) {
          throw new Error(`${error.response?.data?.message || error.message}`);
        } else {
          throw new Error(`${error.message}`);
        }
      } else {
        throw new Error("An unexpected error occurred");
      }
    }
  };

  