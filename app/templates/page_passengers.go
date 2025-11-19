package templates

import (
    "context"
    "io"
    templ "github.com/a-h/templ"
)

func PagePassengers() templ.Component {
    body := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
        _, _ = io.WriteString(w, "<section class=\"mx-auto max-w-6xl p-4\">")
        _, _ = io.WriteString(w, "<div class=\"card bg-base-200/60 border border-white/10 rounded-box shadow-xl ring-1 ring-white/10\">")
        _, _ = io.WriteString(w, "<div class=\"card-body\">")
        _, _ = io.WriteString(w, "<h2 class=\"card-title\">Passengers</h2>")
        _, _ = io.WriteString(w, "<p class=\"opacity-80\">Scaffolded page. Edit at app/templates/page_passengers.go</p>")
        _, _ = io.WriteString(w, "</div></div></section>")
        return nil
    })
    return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error { return LayoutSEO(SEO{Title: "Passengers", Description: "Passengers page", Canonical: "/passengers"}).Render(templ.WithChildren(ctx, body), w) })
}
