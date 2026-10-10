import { Controller } from "@hotwired/stimulus"

// Client-side search over /search.json. Every term must appear in a guide's
// title, section or excerpt. Cmd/Ctrl+K focuses the box; arrows and Enter pick.
export default class extends Controller {
  static targets = [ "input", "results" ]
  static values = { url: String }

  connect() {
    this.shortcut = (event) => {
      if ((event.metaKey || event.ctrlKey) && event.key === "k") {
        event.preventDefault()
        this.inputTarget.focus()
      }
    }
    document.addEventListener("keydown", this.shortcut)
  }

  disconnect() {
    document.removeEventListener("keydown", this.shortcut)
  }

  async query() {
    const terms = this.inputTarget.value.trim().toLowerCase().split(/\s+/).filter(Boolean)
    if (terms.length === 0) return this.hide()

    this.index ||= await fetch(this.urlValue).then((response) => response.json()).catch(() => [])
    const matches = this.index.filter((doc) => {
      const text = `${doc.title} ${doc.section} ${doc.excerpt}`.toLowerCase()
      return terms.every((term) => text.includes(term))
    }).slice(0, 8)

    this.resultsTarget.replaceChildren(...(matches.length ? matches.map((doc) => this.row(doc)) : [ this.empty() ]))
    this.resultsTarget.hidden = false
  }

  keydown(event) {
    const links = [ ...this.resultsTarget.querySelectorAll("a") ]
    if (event.key === "Escape") return this.hide()
    if (!links.length || ![ "ArrowDown", "ArrowUp" ].includes(event.key)) return

    event.preventDefault()
    const current = links.indexOf(document.activeElement)
    const next = event.key === "ArrowDown" ? current + 1 : current - 1
    if (next < 0) this.inputTarget.focus()
    else links[Math.min(next, links.length - 1)].focus()
  }

  hideLater() {
    setTimeout(() => { if (!this.element.contains(document.activeElement)) this.hide() }, 150)
  }

  hide() {
    this.resultsTarget.hidden = true
  }

  row(doc) {
    const link = document.createElement("a")
    link.href = doc.url
    link.className = "search__result"

    const title = document.createElement("span")
    title.className = "search__title"
    title.textContent = doc.title

    const section = document.createElement("span")
    section.className = "search__section"
    section.textContent = doc.section

    link.append(title, section)
    return link
  }

  empty() {
    const note = document.createElement("p")
    note.className = "search__empty"
    note.textContent = "No guides match."
    return note
  }
}
