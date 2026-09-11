import { apiClient } from "@/lib/api/client";

/**
 * Museumに所属するExhibition一覧を取得する。
 *
 * UI Component自身にAPI通信処理を書かず、
 * データ取得責務をこのファイルへ分離する。
 */
export async function getMuseumExhibitions(
  museumId: string,
) {
  const {
    data,
    error,
  } = await apiClient.GET(
    "/museums/{museum_id}/exhibitions",
    {
      params: {
        path: {
          museum_id: museumId,
        },
      },
    },
  );

  if (error || !data) {
    throw new Error(
      "展示室一覧の取得に失敗しました。",
    );
  }

  return data.exhibitions;
}

/**
 * Exhibitionを1件取得する。
 */
export async function getExhibition(
  id: string,
) {
  const {
    data,
    error,
  } = await apiClient.GET(
    "/exhibitions/{id}",
    {
      params: {
        path: {
          id,
        },
      },
    },
  );

  if (error || !data) {
    return null;
  }

  return data;
}