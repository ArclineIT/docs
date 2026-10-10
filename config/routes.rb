Rails.application.routes.draw do
  root "guides#index"

  # Sign-in through Nexus.
  get    "sign_in",             to: "sessions#new",      as: :sign_in
  get    "auth/nexus/callback", to: "sessions#callback", as: :nexus_callback
  delete "logout",              to: "sessions#destroy",  as: :logout

  get "client",       to: "client_pages#index", as: :client_pages
  get "client/:slug", to: "client_pages#show",  as: :client_page

  namespace :admin do
    root to: redirect("/admin/pages")
    resources :pages, except: :show
  end

  get "search.json", to: "feeds#search",  as: :search_index
  get "sitemap.xml", to: "feeds#sitemap", as: :sitemap
  get "rss.xml",     to: "feeds#rss",     as: :rss

  # Health check for load balancers and uptime monitors.
  get "up" => "rails/health#show", as: :rails_health_check

  # Public guides last, so their two-segment paths do not swallow the routes above.
  get ":section",       to: "guides#section", as: :guide_section, constraints: { section: /[a-z0-9-]+/ }
  get ":section/:slug", to: "guides#show",    as: :guide,         constraints: { section: /[a-z0-9-]+/, slug: /[a-z0-9-]+/ }
end
