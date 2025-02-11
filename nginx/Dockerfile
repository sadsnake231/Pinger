FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

RUN chown -R nginx:nginx /usr/share/nginx/html && chmod -R 755 /usr/share/nginx/html


CMD ["nginx", "-g", "daemon off;"]
