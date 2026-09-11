"use client";

import {
  FormEvent,
  useState,
} from "react";
import { useRouter } from "next/navigation";

import { apiClient } from "@/lib/api/client";

/**
 * Exhibition作成フォーム。
 *
 * このComponentは、
 *
 * - Form入力
 * - API送信
 * - 送信状態
 * - エラー表示
 *
 * のみを担当する。
 *
 * Museumページ全体のレイアウトについては
 * page.tsx側へ持たせる。
 */
export default function CreateExhibitionForm() {
  const router =
    useRouter();

  // ----------------------------------------
  // Form state
  // ----------------------------------------

  const [title, setTitle] =
    useState("");

  const [description, setDescription] =
    useState("");

  // API送信中かどうか。
  const [isSubmitting, setIsSubmitting] =
    useState(false);

  // API Error表示用。
  const [errorMessage, setErrorMessage] =
    useState<string | null>(null);

  /**
   * Exhibitionを作成する。
   */
  async function handleSubmit(
    event: FormEvent<HTMLFormElement>,
  ) {
    // Browser標準のForm送信を止める。
    event.preventDefault();

    setIsSubmitting(true);
    setErrorMessage(null);

    try {
      /**
       * URLとRequest Bodyは
       * OpenAPIから自動生成された型によって検証される。
       *
       * museum_idは現時点では仮値。
       *
       * Authentication / Museum取得を実装したら、
       * ログインユーザーのMuseum IDへ置き換える。
       */
      const {
        data,
        error,
      } = await apiClient.POST(
        "/exhibitions",
        {
          body: {
            museum_id: "museum-001",
            title,
            description,
          },
        },
      );

      if (error) {
        setErrorMessage(
          "展示室を作成できませんでした。",
        );

        return;
      }

      if (!data) {
        setErrorMessage(
          "展示室を作成できませんでした。",
        );

        return;
      }

      // 作成成功後、フォームを初期化する。
      setTitle("");
      setDescription("");

      console.log(
        "Exhibition created:",
        data,
      );
    } catch (error) {
      console.error(error);

      setErrorMessage(
        "サーバーとの通信に失敗しました。",
      );
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <form
      className="exhibition-form"
      onSubmit={handleSubmit}
    >
      <div className="exhibition-form__field">
        <label
          htmlFor="exhibition-title"
          className="exhibition-form__label"
        >
          展示室の名前
        </label>

        <input
          id="exhibition-title"
          type="text"
          className="exhibition-form__input"
          placeholder="つくったもの"
          value={title}
          onChange={(event) =>
            setTitle(event.target.value)
          }
          required
        />
      </div>

      <div className="exhibition-form__field">
        <label
          htmlFor="exhibition-description"
          className="exhibition-form__label"
        >
          この展示について
        </label>

        <textarea
          id="exhibition-description"
          className="exhibition-form__textarea"
          placeholder="ここに残しておきたいものについて。"
          value={description}
          onChange={(event) =>
            setDescription(
              event.target.value,
            )
          }
          rows={3}
        />
      </div>

      {errorMessage && (
        <p
          className="exhibition-form__error"
          role="alert"
        >
          {errorMessage}
        </p>
      )}

      <button
        type="submit"
        className="exhibition-form__submit"
        disabled={
          isSubmitting ||
          title.trim() === ""
        }
      >
        {isSubmitting
          ? "展示室をつくっています..."
          : "展示室をつくる"}
      </button>
    </form>
  );
}