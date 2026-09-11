"use client";

import { useState, type FormEvent } from "react";
import { createExhibit, getExhibits } from "@/lib/api/exhibits";
import type { components } from "@/lib/api/generated/schema";

type Props = {
  exhibitionId: string;
  initialExhibits: components["schemas"]["ExhibitResponse"][];
  initialError?: string;
};

export default function ExhibitGallery({ exhibitionId, initialExhibits, initialError }: Props) {
  const [exhibits, setExhibits] = useState(initialExhibits);
  const [listError, setListError] = useState(initialError ?? "");
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [busy, setBusy] = useState(false);

  async function reload() {
    try {
      setExhibits(await getExhibits(exhibitionId));
      setListError("");
    } catch {
      setListError("展示物を取得できませんでした。再読み込みしてください。");
    }
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (busy) return;
    const form = event.currentTarget;
    const values = new FormData(form);
    setBusy(true);
    setError("");
    setNotice("");
    try {
      await createExhibit(exhibitionId, {
        image_url: String(values.get("image_url") ?? "").trim(),
        title: String(values.get("title") ?? ""),
        caption: String(values.get("caption") ?? ""),
      });
      // Saving and reloading have distinct outcomes: a failed GET must not invite
      // another POST for an exhibit that has already been saved.
      form.reset();
      setNotice("展示しました。");
      await reload();
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "展示できませんでした。");
    } finally {
      setBusy(false);
    }
  }

  return (
    <>
      <div className="exhibit-gallery" aria-label="展示物一覧">
        {listError ? (
          <div role="alert">
            <p>{listError}</p>
            <button type="button" className="exhibition-form__submit" onClick={reload}>再読み込み</button>
          </div>
        ) : exhibits.length === 0 ? (
          <p className="exhibition-detail__empty">この展示室には、まだ展示物がありません。</p>
        ) : null}
        {exhibits.map((exhibit) => (
          <figure className="exhibit" key={exhibit.id}>
            {/* External hosts are unrestricted for this MVP; the browser loads the original image. */}
            <picture>
              <img className="exhibit__image" src={exhibit.image_url} alt={exhibit.title || "展示作品"} loading="lazy" referrerPolicy="no-referrer" />
            </picture>
            <figcaption className="exhibit__label">
              {exhibit.title && <h2 className="exhibit__title">{exhibit.title}</h2>}
              {exhibit.caption && <p className="exhibit__caption">{exhibit.caption}</p>}
            </figcaption>
          </figure>
        ))}
      </div>
      <section className="create-exhibit" aria-labelledby="create-exhibit-heading">
        <h2 id="create-exhibit-heading" className="section-header__title">新しい展示</h2>
        <form className="exhibition-form" onSubmit={submit}>
          <div className="exhibition-form__field">
            <label className="exhibition-form__label" htmlFor="exhibit-image">画像URL</label>
            <input id="exhibit-image" name="image_url" type="url" required pattern="https?://.+" className="exhibition-form__input exhibit-form__url" placeholder="https://example.com/image.jpg" />
          </div>
          <div className="exhibition-form__field">
            <label className="exhibition-form__label" htmlFor="exhibit-title">作品タイトル（任意）</label>
            <input id="exhibit-title" name="title" className="exhibition-form__input" />
          </div>
          <div className="exhibition-form__field">
            <label className="exhibition-form__label" htmlFor="exhibit-caption">キャプション（任意）</label>
            <textarea id="exhibit-caption" name="caption" className="exhibition-form__textarea" rows={3} />
          </div>
          {error && <p className="exhibition-form__error" role="alert">{error}</p>}
          <p role="status">{notice}</p>
          <button className="exhibition-form__submit" disabled={busy}>{busy ? "展示しています…" : "展示する"}</button>
        </form>
      </section>
    </>
  );
}
