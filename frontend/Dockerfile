FROM node:alpine as build
ARG PUBLIC_FRONTEND_URL
ARG PRIVATE_ALLOW_SEARCH_INDEXING
WORKDIR /app
COPY . .
RUN npm i -g pnpm && \
    pnpm i &&\
    if [[ -n "$PUBLIC_FRONTEND_URL" && "$PRIVATE_ALLOW_SEARCH_INDEXING" == "true" ]]; then echo -e "User-agent: *\nDisallow: /login\nDisallow: /admin\nDisallow: /logout\n\nSitemap: ${PUBLIC_FRONTEND_URL}/sitemap.xml" > static/robots.txt; else echo -e "User-agent: *\nDisallow: /" > static/robots.txt; fi && \
    pnpm run build

FROM node:alpine
USER node:node
WORKDIR /app
COPY --from=build --chown=node:node /app/build ./build
COPY --from=build --chown=node:node /app/node_modules ./node_modules
COPY --chown=node:node package.json .
CMD ["node","build"]
