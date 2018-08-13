package xtraefik

func Tags() []string {
	return []string{"traefik.enable=true"}
}

func TagsForGrpc() []string {
	tags := Tags()
	tags = append(tags, "traefik.protocol=h2c")
	return tags
}
