# Instapod

Instapod is as simple server that lets you add listen to YouTube videos in a podcast player. It uses youtube-dl to fetch the video and convert it to audio, and then serves the files as a podcast feed.

On first run, the server will generate a config file. Update the `Title` and `Description` fields to your liking and `BaseURL` to point to your server. You should also set `Link` and `ImageURL` so that you have a link in your metadata and nice album art for your feed.