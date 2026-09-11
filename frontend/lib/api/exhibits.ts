import { apiClient } from "./client";
import type { components } from "./generated/schema";

export async function getExhibits(exhibitionId: string) {
  const { data, error } = await apiClient.GET("/exhibitions/{exhibition_id}/exhibits", {
    params: { path: { exhibition_id: exhibitionId } },
    cache: "no-store",
  });
  if (error || !data) throw new Error("展示物を取得できませんでした。");
  return data.exhibits;
}

export async function createExhibit(
  exhibitionId: string,
  input: components["schemas"]["CreateExhibitRequest"],
) {
  const { data, error } = await apiClient.POST("/exhibitions/{exhibition_id}/exhibits", {
    params: { path: { exhibition_id: exhibitionId } },
    body: input,
  });
  if (error || !data) throw new Error("展示できませんでした。入力内容と展示室を確認してください。");
  return data;
}
