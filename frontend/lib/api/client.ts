import createClient from "openapi-fetch";

import type { paths } from "./generated/schema";

/**
 * Backend APIへアクセスするための共通Client。
 *
 * 各Componentで直接fetch()を書かず、
 * 必ずこのClientを経由する。
 *
 * OpenAPIから生成されたpaths型を利用することで、
 *
 * - URL
 * - Request
 * - Response
 * - HTTP Method
 *
 * をTypeScriptで検証できる。
 */
export const apiClient = createClient<paths>({
  /**
   * Local Backend API。
   *
   * 後ほど環境変数へ移動するが、
   * 最初の動作確認ではlocalhostを使用する。
   */
  baseUrl: "http://localhost:8080",
});