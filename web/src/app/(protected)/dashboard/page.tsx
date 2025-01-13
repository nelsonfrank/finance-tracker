import { userDashboardAPI } from "@/data/backend/api/dashboard";
// import { useEffect } from "react";

export default async function DashboardPage() {
  // useEffect(() => {
  //   fetchData();
  // }, []);

  // async function fetchData() {
  try {
    const data = await userDashboardAPI();
    console.log(data.data);
  } catch (error) {
    console.log({ error });
  }
  // }

  return <p>Protected route</p>;
}
