"use client"
import { useEffect } from "react";
import { userDashboardAPI } from "@/data/backend/api/dashboard";
import { Button } from "@/components/ui/button";

export default function DashboardPage() {
  useEffect(() => {
    fetchData();
  }, []);

  async function fetchData() {
  try {
    const data = await userDashboardAPI();
    console.log(data);
  } catch (error) {
    console.log({ error });
  }
  }

  return <div>
    <Button onClick={fetchData}>Fetch Data</Button>
  </div>;
}
