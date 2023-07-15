# TODO

- [ ] Allow environment variables to dictate data sources at runtime.
- [ ] split site into olympiad specific and short answer specific

- [ ] add some kind of UI for people to drag and drop problems of their choice
- [ ] dark mode

- [x] change data.json layout to include the name of the contests that were scraped
- [x] change input json to be more flexible and accept single contests instead of just collections of contests

- [ ] implement endpoint for users to see individual rendered problems.

# Proposals

- [x] Move to Meilisearch
- [x] Isolate db query logic to backend/database
- [ ] Change site architecture and move to SSR (So that problems get
      pre-rendered)
    - [ ] Configure KaTeX so that it's not client-side dependent.
