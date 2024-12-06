FROM alpine:3.6
 
RUN apk update && \
    apk add --no-cache nginx
 
ADD src /app
ADD ./src/nginx.conf /etc/nginx/conf.d/default.conf
 
RUN mkdir -p /run/nginx
 
CMD ["nginx", "-g", "daemon off;"]
