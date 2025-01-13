"use client";

import { useEffect } from "react";

interface GlobalErrorProps {
  error: Error; // Represents the error object
  reset: () => void; // Function to reset the error boundary
}
export default function GlobalError({ error, reset }: GlobalErrorProps) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <html>
      <body>
        <h1>Something went wrong!</h1>
        <p>{error.message}</p>
        <button onClick={() => reset()}>Try again</button>
      </body>
    </html>
  );
}
