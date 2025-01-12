import { type NextRequest } from "next/server";
import * as API from "@/data/backend/api";

export async function POST(req: NextRequest) {
  const payload = await req.json();
  const response = await API.registerUserAPI(payload);

  const data = response.data;

  return Response.json(data);
}
