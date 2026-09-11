import ExhibitGallery from "@/components/exhibit/ExhibitGallery";
import { getExhibits } from "@/lib/api/exhibits";

import Link from "next/link";
import { notFound } from "next/navigation";

import {
  getExhibition,
} from "@/lib/api/exhibitions";

type ExhibitionPageProps = {
  params: Promise<{
    id: string;
  }>;
};

/**
 * Exhibition詳細ページ。
 *
 * URL:
 *
 * /exhibitions/{id}
 *
 * Backend:
 *
 * GET /exhibitions/{id}
 */
export default async function ExhibitionPage(
  props: ExhibitionPageProps,
) {
  const { id } =
    await props.params;

  const exhibition =
    await getExhibition(id);

  if (!exhibition) {
    notFound();
  }

  const gallery = await getExhibits(id)
    .then((exhibits) => ({ exhibits, error: undefined }))
    .catch(() => ({ exhibits: [], error: "展示物を取得できませんでした。" }));

  return (
    <main className="museum">
      <header className="museum-header">
        <div className="museum-header__identity">
          <Link href="/">
            ← ミュージアムへ戻る
          </Link>
        </div>
      </header>

      <section className="exhibition-detail">
        <p className="exhibition-detail__label">
          展示室
        </p>

        <h1 className="exhibition-detail__title">
          {exhibition.title}
        </h1>

        {exhibition.description && (
          <p className="exhibition-detail__description">
            {
              exhibition.description
            }
          </p>
        )}

        <ExhibitGallery key={id} exhibitionId={id} initialExhibits={gallery.exhibits} initialError={gallery.error} />
      </section>
    </main>
  );
}