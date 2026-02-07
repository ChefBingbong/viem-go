import { NextResponse } from "next/server";
import { buildSearchIndex } from "@/lib/search";

export const dynamic = "force-static";

export async function GET() {
  const index = buildSearchIndex();
  return NextResponse.json(index);
}
