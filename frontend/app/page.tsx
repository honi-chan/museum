import Link from "next/link";

import CreateExhibitionForm from "@/components/exhibition/CreateExhibitionForm";
import {
  getMuseumExhibitions,
} from "@/lib/api/exhibitions";

/**
 * 現時点ではAuthenticationがないため、
 * 開発用Museum IDを固定している。
 *
 * Authentication実装後は
 * Login UserのMuseum IDへ置き換える。
 */
const MUSEUM_ID =
  "museum-001";

/**
 * Museum Home
 *
 * Server Componentとして、
 * BackendからExhibition一覧を取得して表示する。
 */
export default async function Home() {
  const exhibitions =
    await getMuseumExhibitions(
      MUSEUM_ID,
    );

  return (
    <main className="museum">
      <header className="museum-header">
        <div className="museum-header__identity">
          <p className="museum-header__eyebrow">
            わたしのミュージアム
          </p>

          <h1 className="museum-header__title">
            Haruki&apos;s Museum
          </h1>
        </div>

        <nav
          className="museum-header__navigation"
          aria-label="ミュージアムナビゲーション"
        >
          <a href="#exhibitions">
            展示室
          </a>
        </nav>
      </header>

      <section className="museum-introduction">
        <p className="museum-introduction__number">
          MUSEUM 001
        </p>

        <h2 className="museum-introduction__statement">
          つくったもの。
          <br />
          好きだったもの。
          <br />
          忘れたくないもの。
        </h2>

        <p className="museum-introduction__description">
          ここは、わたしという人間を
          少しずつ残していく場所です。
          大切にしたいものだけを、
          展示しています。
        </p>
      </section>

      <section
        id="exhibitions"
        className="exhibitions"
      >
        <header className="section-header">
          <p className="section-header__number">
            01
          </p>

          <h2 className="section-header__title">
            展示室
          </h2>
        </header>

        <div className="exhibition-list">
          {exhibitions.map(
            (
              exhibition,
              index,
            ) => (
              <Link
                key={exhibition.id}
                href={`/exhibitions/${exhibition.id}`}
                className="exhibition"
              >
                <p className="exhibition__number">
                  {String(
                    index + 1,
                  ).padStart(
                    3,
                    "0",
                  )}
                </p>

                <div>
                  <h3 className="exhibition__title">
                    {
                      exhibition.title
                    }
                  </h3>

                  <p className="exhibition__description">
                    {
                      exhibition.description
                    }
                  </p>
                </div>
              </Link>
            ),
          )}
        </div>
      </section>

      <section className="create-exhibition">
        <header className="section-header">
          <p className="section-header__number">
            02
          </p>

          <h2 className="section-header__title">
            新しい展示室
          </h2>
        </header>

        <div className="create-exhibition__content">
          <p className="create-exhibition__description">
            残しておきたいもののために、
            新しい展示室をつくります。
          </p>

          <CreateExhibitionForm />
        </div>
      </section>
    </main>
  );
}