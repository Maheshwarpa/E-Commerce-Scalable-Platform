FROM docker.elastic.co/logstash/logstash:7.17.5

# Copy Logstash pipeline configuration
COPY pipeline/logstash.conf /usr/share/logstash/pipeline/logstash.conf

CMD ["logstash", "-f", "/usr/share/logstash/pipeline/logstash.conf"]
