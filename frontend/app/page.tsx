import CreateExhibitionForm from "@/components/exhibition/CreateExhibitionForm";

/**
 * Museum Home
 *
 * ユーザー自身のMuseumを表すトップページ。
 *
 * 一般的なSNSのDashboardではなく、
 * 美術館の入口・展示案内のような構成を目指す。
 *
 * このpage.tsx自身はServer Component。
 *
 * Interactionが必要な部分だけ
 * Client Componentとして分離する。
 */
export default function Home() {
  return (
    <main className="museum">
      {/* ------------------------------------
       * Museum Header
       * ------------------------------------ */}
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

          <a href="#about">
            この場所について
          </a>
        </nav>
      </header>

      {/* ------------------------------------
       * Introduction
       * ------------------------------------ */}
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

      {/* ------------------------------------
       * Exhibitions
       * ------------------------------------ */}
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
          <article className="exhibition">
            <p className="exhibition__number">
              001
            </p>

            <div>
              <h3 className="exhibition__title">
                つくったもの
              </h3>

              <p className="exhibition__description">
                作品、実験、アイデア。
                自分の手から生まれたものたち。
              </p>
            </div>
          </article>

          <article className="exhibition">
            <p className="exhibition__number">
              002
            </p>

            <div>
              <h3 className="exhibition__title">
                京都の夏
              </h3>

              <p className="exhibition__description">
                忘れたくなかった、
                あの夏の断片。
              </p>
            </div>
          </article>

          <article className="exhibition">
            <p className="exhibition__number">
              003
            </p>

            <div>
              <h3 className="exhibition__title">
                好きなもの
              </h3>

              <p className="exhibition__description">
                もの、物語、場所。
                ずっと心に残っているものたち。
              </p>
            </div>
          </article>
        </div>
      </section>

      {/* ------------------------------------
       * Exhibition Creation
       * ------------------------------------ */}
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