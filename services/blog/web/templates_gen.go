package web

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _templateAssetsa3d57443642ccfe6165226140fc22a3f996d43ef = "{{define \"title\"}}サインイン{{end}}\n\n{{define \"body\"}}\n<section>\n  <form method=\"POST\" action=\"/signin\">\n    <input type=\"hidden\" name=\"csrf_token\" value=\"{{.CsrfToken}}\">\n    <p><label>ユーザー名 <input type=\"text\" name=\"name\"></label></p>\n    <p><label>パスワード <input type=\"text\" name=\"password\"></label></p>\n    <p><input type=\"submit\" value=\"サインイン\"></p>\n  </form>\n</section>\n{{end}}\n"
var _templateAssets8a236120b99e2d6e5cb52991b22412edafc7a6cf = "{{define \"title\"}}{{.Blog.Title}} を編集{{end}}\n\n{{define \"body\"}}\n<section>\n  <header>\n    <h1>{{.Blog.Title}} を編集</h1>\n    <p><a href=\"/my/blogs/{{.Blog.Path}}\">ブログ詳細に戻る</a></p>\n  </header>\n  <section>\n    <form method=\"POST\" action=\"/my/blogs/{{.Blog.Path}}/edit\">\n      <input type=\"hidden\" name=\"csrf_token\" value=\"{{.CsrfToken}}\">\n      <p><label>パス <span>/blogs/{{.Blog.Path}}</span></p>\n      <p><label>タイトル <input type=\"text\" name=\"title\" value=\"{{.Blog.Title}}\"></label></p>\n      <p><label>説明 <textarea name=\"description\">{{.Blog.Description}}</textarea></label></p>\n      <p><input type=\"submit\" value=\"更新\"></p>\n    </form>\n  </section>\n  <section>\n    <form method=\"POST\" action=\"/my/blogs/{{.Blog.Path}}/delete\">\n      <input type=\"hidden\" name=\"csrf_token\" value=\"{{.CsrfToken}}\">\n      <p><input type=\"submit\" value=\"削除\"></p>\n    </form>\n  </section>\n</section>\n{{end}}\n"
var _templateAssets782418f0d2f8f57038b0d540ac47cea598e7eb2e = "{{define \"title\"}}エントリを編集{{end}}\n\n{{define \"body\"}}\n<section>\n  <header>\n    <h1>エントリを編集</h1>\n    <p><a href=\"/my/blogs/{{.Blog.Path}}\">ブログ詳細に戻る</a></p>\n  </header>\n  <section>\n    <form method=\"POST\" action=\"/my/blogs/{{.Blog.Path}}/entries/{{.Entry.ID}}/edit\">\n      <input type=\"hidden\" name=\"csrf_token\" value=\"{{.CsrfToken}}\">\n      <p><label>タイトル <input type=\"text\" name=\"title\" value=\"{{.Entry.Title}}\"></label></p>\n      <p><label>本文 <textarea name=\"body\">{{.Entry.Body}}</textarea></label></p>\n      <p><input type=\"submit\" value=\"更新\"></p>\n    </form>\n  </section>\n  <section>\n    <form method=\"POST\" action=\"/my/blogs/{{.Blog.Path}}/entries/{{.Entry.ID}}/unpublish\">\n      <input type=\"hidden\" name=\"csrf_token\" value=\"{{.CsrfToken}}\">\n      <p><input type=\"submit\" value=\"削除\"></p>\n    </form>\n  </section>\n</section>\n{{end}}\n"
var _templateAssets8cc4c7c9d609cc42b871d5e1625b5fb9eb83913e = "{{define \"title\"}}新規エントリ投稿{{end}}\n\n{{define \"body\"}}\n<section>\n  <header>\n    <h1>新規エントリ投稿</h1>\n    <p><a href=\"/my/blogs/{{.Blog.Path}}\">ブログ詳細に戻る</a></p>\n  </header>\n  <section>\n    <form method=\"POST\" action=\"/my/blogs/{{.Blog.Path}}/entries/-/publish\">\n      <input type=\"hidden\" name=\"csrf_token\" value=\"{{.CsrfToken}}\">\n      <p><label>タイトル <input type=\"text\" name=\"title\"></label></p>\n      <p><label>本文 <textarea name=\"body\"></textarea></label></p>\n      <p><input type=\"submit\" value=\"投稿\"></p>\n    </form>\n  </section>\n</section>\n{{end}}\n"
var _templateAssets4423b5123543a5e077f0ec2952bf061f105f3950 = "{{define \"title\"}}サインアウト{{end}}\n\n{{define \"body\"}}\n<section>\n  <form method=\"POST\" action=\"/signout\">\n    <input type=\"hidden\" name=\"csrf_token\" value=\"{{.CsrfToken}}\">\n    <p><input type=\"submit\" value=\"サインアウト\"></p>\n  </form>\n</section>\n{{end}}\n"
var _templateAssetsfaeaa66ba924cb7d86e0551ff94c8bba26df5b30 = "{{define \"title\"}}{{.Blog.Title}}{{end}}\n\n{{define \"body\"}}\n<section>\n  <header>\n    <h1>{{.Blog.Title}}</h1>\n    <p><a href=\"/my/blogs\">ブログ一覧に戻る</a></p>\n    <p><a href=\"/blogs/{{.Blog.Path}}\" target=\"_blank\" rel=\"nofollow noopener\">Go</a></p>\n  </header>\n  <section>\n    <dl>\n      <dt>パス</dt>\n      <dd>/blogs/{{.Blog.Path}}</dd>\n      <dt>説明</dt>\n      <dd>{{.Blog.Description}}</dd>\n    </dl>\n    <p><a href=\"/my/blogs/{{.Blog.Path}}/edit\">編集</a></p>\n  </section>\n  <section>\n    <h1>エントリ一覧</h1>\n    <p><a href=\"/my/blogs/{{.Blog.Path}}/entries/-/publish\">新規投稿</a></p>\n    <table>\n      <thead>\n        <tr>\n          <td>タイトル</td>\n          <td>投稿日時</td>\n          <td>編集日時</td>\n          <td></td>\n        </tr>\n      </thead>\n      <tbody>\n        {{range .Entries}}\n        <tr>\n          <td><a href=\"/my/blogs/{{$.Blog.Path}}/entries/{{.ID}}\">{{.Title}}</a></td>\n          <td>{{.PublishedAt}}</td>\n          <td>{{.EditedAt}}</td>\n          <td><a href=\"/blogs/{{$.Blog.Path}}/entries/{{.ID}}\" target=\"_blank\" rel=\"nofollow noopener\">Go</a></td>\n        </tr>\n        {{end}}\n      </tbody>\n    </table>\n    <footer>\n      <p>\n        <span>\n          {{if .HasPrevPage}}\n          <a href=\"/my/blogs/{{.Blog.Path}}?page={{.PrevPage}}\">&lt;</a>\n          {{end}}\n        </span>\n        <span>\n          Page {{.Page}}\n        </span>\n        <span>\n          {{if .HasNextPage}}\n          <a href=\"/my/blogs/{{.Blog.Path}}?page={{.NextPage}}\">&gt;</a>\n          {{end}}\n        </span>\n      </p>\n    </footer>\n  </section>\n</section>\n{{end}}\n"
var _templateAssetsb7628a386984abc0fc223cfeb88f251d466ad8e9 = "<!DOCTYPE html>\n<html>\n  <head>\n    <meta charset=\"UTF-8\">\n    <title>{{block \"title\" .}}{{end}}</title>\n    {{block \"head\" .}}{{end}}\n  </head>\n  <body>\n    {{block \"body\" .}}{{end}}\n  </body>\n</html>\n"
var _templateAssets9f0ae851af9adb08e3242917d7d52dc270c7bae2 = "{{define \"title\"}}ブログ{{end}}\n\n{{define \"body\"}}\n<section>\n  <header>\n    <h1>ブログ</h1>\n    <nav>\n      {{if .User}}\n      <p>ようこそ {{.User.Name}} さん</p>\n      {{end}}\n      <ul>\n        {{if .User}}\n        <li><a href=\"/my/blogs\">マイページ</a></li>\n        <li><a href=\"/signout\">サインアウト</a></li>\n        {{else}}\n        <li><a href=\"/signup\">サインアップ</a></li>\n        <li><a href=\"/signin\">サインイン</a></li>\n        {{end}}\n      </ul>\n    </nav>\n  </header>\n  <section>\n    <header>\n      <h1>みんなのブログ</h1>\n    </header>\n    <table>\n      <thead>\n        <tr>\n          <td>タイトル</td>\n          <td>説明</td>\n        </tr>\n      </thead>\n      <tbody>\n        {{range .Blogs}}\n        <tr>\n          <td><a href=\"/blogs/{{.Path}}\">{{.Title}}</a></td>\n          <td>{{.Description}}</td>\n        </tr>\n        {{end}}\n      </tbody>\n    </table>\n    <footer>\n      <p>\n        <span>\n          {{if .HasPrevPage}}\n          <a href=\"/?page={{.PrevPage}}\">&lt;</a>\n          {{end}}\n        </span>\n        <span>\n          Page {{.Page}}\n        </span>\n        <span>\n          {{if .HasNextPage}}\n          <a href=\"/?page={{.NextPage}}\">&gt;</a>\n          {{end}}\n        </span>\n      </p>\n    </footer>\n  </section>\n</section>\n{{end}}\n"
var _templateAssets56684805376519c6927b11bff43b6315dc39cfc8 = "{{define \"title\"}}{{.Entry.Title}}{{end}}\n\n{{define \"body\"}}\n<section>\n  <header>\n    <h1><a href=\"/blogs/{{.Blog.Path}}\">{{.Blog.Title}}</a></h1>\n    <p>{{.Blog.Description}}</p>\n    {{if .IsAuthor}}\n    <p><a href=\"/my/blogs/{{.Blog.Path}}\">管理</a></p>\n    {{end}}\n  </header>\n  <section>\n    <article>\n      <header>\n        <h1><a href=\"/blogs/{{.Blog.Path}}/entries/{{.Entry.ID}}\">{{.Entry.Title}}</a></h1>\n        <p>{{.Entry.PublishedAt}}</p>\n        {{if .IsAuthor}}\n        <p><a href=\"/my/blogs/{{.Blog.Path}}/entries/{{.Entry.ID}}\">編集</a></p>\n        {{end}}\n      </header>\n      <section>{{.Entry.BodyHTML | unescapedHTML}}</section>\n    </article>\n  </section>\n</section>\n{{end}}\n"
var _templateAssets219ff4155ea7d695814c4519a438e5b1de933def = "{{define \"title\"}}{{.Blog.Title}}{{end}}\n\n{{define \"body\"}}\n<section>\n  <header>\n    <h1><a href=\"/blogs/{{.Blog.Path}}\">{{.Blog.Title}}</a></h1>\n    <p>{{.Blog.Description}}</p>\n    {{if .IsAuthor}}\n    <p><a href=\"/my/blogs/{{.Blog.Path}}\">管理</a></p>\n    {{end}}\n  </header>\n  <section>\n    {{range .Entries}}\n    <article>\n      <header>\n        <h1><a href=\"/blogs/{{$.Blog.Path}}/entries/{{.ID}}\">{{.Title}}</a></h1>\n        <p>{{.PublishedAt}}</p>\n        {{if $.IsAuthor}}\n        <p><a href=\"/my/blogs/{{$.Blog.Path}}/entries/{{.ID}}\">編集</a></p>\n        {{end}}\n      </header>\n      <section>{{.BodyHTML | unescapedHTML}}</section>\n    </article>\n    {{end}}\n    <footer>\n      <p>\n        <span>\n          {{if .HasPrevPage}}\n          <a href=\"/blogs/{{.Blog.Path}}?page={{.PrevPage}}\">&lt;</a>\n          {{end}}\n        </span>\n        <span>\n          Page {{.Page}}\n        </span>\n        <span>\n          {{if .HasNextPage}}\n          <a href=\"/blogs/{{.Blog.Path}}?page={{.NextPage}}\">&gt;</a>\n          {{end}}\n        </span>\n      </p>\n    </footer>\n  </section>\n</section>\n{{end}}\n"
var _templateAssets31fe538b2a987132914302d14e1235f1a59e9d13 = "{{define \"title\"}}{{.User.Name}} のブログ一覧{{end}}\n\n{{define \"body\"}}\n<section>\n  <header>\n    <h1>{{.User.Name}} のブログ一覧</h1>\n    <p><a href=\"/my/blogs/-/create\">新規作成</a></p>\n  </header>\n  <section>\n    <table>\n      <thead>\n        <tr>\n          <td>タイトル</td>\n          <td>説明</td>\n          <td></td>\n        </tr>\n      </thead>\n      <tbody>\n        {{range .Blogs}}\n        <tr>\n          <td><a href=\"/my/blogs/{{.Path}}\">{{.Title}}</a></td>\n          <td>{{.Description}}</td>\n          <td><a href=\"/blogs/{{.Path}}\" target=\"_blank\" rel=\"nofollow noopener\">Go</a></td>\n        </tr>\n        {{end}}\n      </tbody>\n    </table>\n    <footer>\n      <p>\n        <span>\n          {{if .HasPrevPage}}\n          <a href=\"/my/blogs?page={{.PrevPage}}\">&lt;</a>\n          {{end}}\n        </span>\n        <span>\n          Page {{.Page}}\n        </span>\n        <span>\n          {{if .HasNextPage}}\n          <a href=\"/my/blogs?page={{.NextPage}}\">&gt;</a>\n          {{end}}\n        </span>\n      </p>\n    </footer>\n  </section>\n</section>\n{{end}}\n"
var _templateAssets751a02880e76c99ff92ae4612d6ec3e1bbe8405a = "{{define \"title\"}}新規ブログ作成{{end}}\n\n{{define \"body\"}}\n<section>\n  <header>\n    <h1>新規ブログ作成</h1>\n    <p><a href=\"/my/blogs\">ブログ一覧に戻る</a></p>\n  </header>\n  <section>\n    <form method=\"POST\" action=\"/my/blogs\">\n      <input type=\"hidden\" name=\"csrf_token\" value=\"{{.CsrfToken}}\">\n      <p><label>パス <input type=\"text\" name=\"path\"></label></p>\n      <p><label>タイトル <input type=\"text\" name=\"title\"></label></p>\n      <p><label>説明 <textarea name=\"description\"></textarea></label></p>\n      <p><input type=\"submit\" value=\"作成\"></p>\n    </form>\n  </section>\n</section>\n{{end}}\n"
var _templateAssets08426c572e59ef1eb67dcf68e66e44254806cfab = "{{define \"title\"}}サインアップ{{end}}\n\n{{define \"body\"}}\n<section>\n  <form method=\"POST\" action=\"/signup\">\n    <input type=\"hidden\" name=\"csrf_token\" value=\"{{.CsrfToken}}\">\n    <p><label>ユーザー名 <input type=\"text\" name=\"name\"></label></p>\n    <p><label>パスワード <input type=\"text\" name=\"password\"></label></p>\n    <p><input type=\"submit\" value=\"サインアップ\"></p>\n  </form>\n</section>\n{{end}}\n"

// templateAssets returns go-assets FileSystem
var templateAssets = assets.NewFileSystem(map[string][]string{}, map[string]*assets.File{
	"signout.html": &assets.File{
		Path:     "signout.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594609595, 1594609595769827680),
		Data:     []byte(_templateAssets4423b5123543a5e077f0ec2952bf061f105f3950),
	}, "my-blog.html": &assets.File{
		Path:     "my-blog.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594609754, 1594609754620985603),
		Data:     []byte(_templateAssetsfaeaa66ba924cb7d86e0551ff94c8bba26df5b30),
	}, "wrapper.html": &assets.File{
		Path:     "wrapper.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594626993, 1594626993627733685),
		Data:     []byte(_templateAssetsb7628a386984abc0fc223cfeb88f251d466ad8e9),
	}, "index.html": &assets.File{
		Path:     "index.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594863393, 1594863393409596785),
		Data:     []byte(_templateAssets9f0ae851af9adb08e3242917d7d52dc270c7bae2),
	}, "entry.html": &assets.File{
		Path:     "entry.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594625535, 1594625535893898035),
		Data:     []byte(_templateAssets56684805376519c6927b11bff43b6315dc39cfc8),
	}, "blog.html": &assets.File{
		Path:     "blog.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594625532, 1594625532750502657),
		Data:     []byte(_templateAssets219ff4155ea7d695814c4519a438e5b1de933def),
	}, "my-blogs.html": &assets.File{
		Path:     "my-blogs.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594609765, 1594609765841506340),
		Data:     []byte(_templateAssets31fe538b2a987132914302d14e1235f1a59e9d13),
	}, "my-blogs-create.html": &assets.File{
		Path:     "my-blogs-create.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594609709, 1594609709339046725),
		Data:     []byte(_templateAssets751a02880e76c99ff92ae4612d6ec3e1bbe8405a),
	}, "signup.html": &assets.File{
		Path:     "signup.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594609612, 1594609612565612815),
		Data:     []byte(_templateAssets08426c572e59ef1eb67dcf68e66e44254806cfab),
	}, "my-blog-edit.html": &assets.File{
		Path:     "my-blog-edit.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594609671, 1594609671249498354),
		Data:     []byte(_templateAssets8a236120b99e2d6e5cb52991b22412edafc7a6cf),
	}, "my-entry.html": &assets.File{
		Path:     "my-entry.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594609723, 1594609723594687780),
		Data:     []byte(_templateAssets782418f0d2f8f57038b0d540ac47cea598e7eb2e),
	}, "my-entries-publish.html": &assets.File{
		Path:     "my-entries-publish.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594624941, 1594624941101808136),
		Data:     []byte(_templateAssets8cc4c7c9d609cc42b871d5e1625b5fb9eb83913e),
	}, "signin.html": &assets.File{
		Path:     "signin.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1594609582, 1594609582516761847),
		Data:     []byte(_templateAssetsa3d57443642ccfe6165226140fc22a3f996d43ef),
	}}, "")
