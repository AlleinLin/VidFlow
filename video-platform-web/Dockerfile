FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:20-alpine AS runner

WORKDIR /app

RUN addgroup -g 1001 -S nodejs && \
    adduser -S sveltekit -u 1001

COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules

USER sveltekit

EXPOSE 3000

ENV NODE_ENV=production
ENV PORT=3000
ENV HOST=0.0.0.0

CMD ["node", "build"]
