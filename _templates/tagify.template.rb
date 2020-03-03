class Tagify < Formula
  desc "Produces a set of tags from given source. Source can be either an HTML page, Markdown document or a plain text. Support English and Russian words."
  homepage "https://www.zoomio.org/tagify"
  url "https://github.com/zoomio/tagify/archive/${VERSION}.tar.gz"
  sha256 "${SHA}"
  
  depends_on "go" => :build
  
  def install
    system "env", "GOOS=darwin", "GOARCH=amd64", "go", "build", *std_go_args, "cmd/cli/cli.go"
  end
  
  test do
    system "go", "test", "./..."
  end
end
