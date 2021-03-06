package ebookdownloader

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Aiicy/htmlquery"
	"github.com/schollz/progressbar/v2"
)

// https://www.999xs.com/

/*
bookid 规则



 https://www.999xs.com/files/article/html/0/591/ -> bookid = 591
 591 -> {0,591}

 https://www.999xs.com/files/article/html/1/1599/ -> bookid = 1599
 1599 -> {1,599}

https://www.999xs.com/files/article/html/75/75842/ -> bookid = 75842
 75842 - > {75,842}


 https://www.999xs.com/files/article/html/113/113582/ -> bookid = 113582
 113582 -> {113,582}
*/

var _ EBookDLInterface = Ebook999XS{}

//999小说网 999xs.com
type Ebook999XS struct {
	URL  string
	Lock *sync.Mutex
}

func New999XS() Ebook999XS {
	return Ebook999XS{
		URL:  "https://www.999xs.com",
		Lock: new(sync.Mutex),
	}
}

/*
input: 113582 ,output: 113/113582
input: 75842, output: 75/75842
input: 1599, output: 1/1599
input: 591, output: 0/591
*/
func handleBookid(bookid string) string {
	tmp := []rune(bookid)
	if len(tmp) == 3 {
		return "0/" + string(tmp)
	}
	//最后结果不算怎样，都保留最后三个数字，前面n个数字需要分离出来
	return (string(tmp[:len(tmp)-3]) + "/" + string(tmp))
}

//GetBookBriefInfo 获取小说的信息
func (this Ebook999XS) GetBookBriefInfo(bookid string, proxy string) BookInfo {

	var bi BookInfo
	pollURL := this.URL + "/" + "files/article/html/" + handleBookid(bookid) + "/"

	//当 proxy 不为空的时候，表示设置代理
	if proxy != "" {
		doc, err := htmlquery.LoadURLWithProxy(pollURL, proxy)
		if err != nil {
			fmt.Println(err.Error())
		}

		//获取书名字
		bookNameMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:book_name']")
		bookName := htmlquery.SelectAttr(bookNameMeta, "content")
		fmt.Println("书名 = ", bookName)

		//获取书作者
		AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
		author := htmlquery.SelectAttr(AuthorMeta, "content")
		fmt.Println("作者 = ", author)

		//获取书的描述信息
		DescriptionMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:description']")
		description := htmlquery.SelectAttr(DescriptionMeta, "content")
		fmt.Println("简介 = ", description)

		//导入信息
		bi = BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
		}
	} else { //没有设置代理
		doc, err := htmlquery.LoadURL(pollURL)
		if err != nil {
			fmt.Println(err.Error())
		}

		//获取书名字
		bookNameMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:book_name']")
		bookName := htmlquery.SelectAttr(bookNameMeta, "content")
		fmt.Println("书名 = ", bookName)

		//获取书作者
		AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
		author := htmlquery.SelectAttr(AuthorMeta, "content")
		fmt.Println("作者 = ", author)

		//获取书的描述信息
		DescriptionMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:description']")
		description := htmlquery.SelectAttr(DescriptionMeta, "content")
		fmt.Println("简介 = ", description)

		//导入信息
		bi = BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
		}
	}
	return bi
}

func (this Ebook999XS) GetBookInfo(bookid string, proxy string) BookInfo {

	var bi BookInfo
	var chapters []Chapter
	pollURL := this.URL + "/" + "files/article/html/" + handleBookid(bookid) + "/"

	//当 proxy 不为空的时候，表示设置代理
	if proxy != "" {
		doc, err := htmlquery.LoadURLWithProxy(pollURL, proxy)
		if err != nil {
			fmt.Println(err.Error())
		}

		//获取书名字
		bookNameMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:book_name']")
		bookName := htmlquery.SelectAttr(bookNameMeta, "content")
		fmt.Println("书名 = ", bookName)

		//获取书作者
		AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
		author := htmlquery.SelectAttr(AuthorMeta, "content")
		fmt.Println("作者 = ", author)

		//获取书的描述信息
		DescriptionMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:description']")
		description := htmlquery.SelectAttr(DescriptionMeta, "content")
		fmt.Println("简介 = ", description)

		//获取书章节列表
		ddNode, _ := htmlquery.Find(doc, "//div[@id='list']//dl//dd")
		for i := 0; i < len(ddNode); i++ {
			var tmp Chapter
			aNode, _ := htmlquery.Find(ddNode[i], "//a")
			tmp.Link = this.URL + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			chapters = append(chapters, tmp)
		}

		//导入信息
		bi = BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
			Chapters:    chapters,
		}
	} else { //没有设置代理
		doc, err := htmlquery.LoadURL(pollURL)
		if err != nil {
			fmt.Println(err.Error())
		}

		//获取书名字
		bookNameMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:book_name']")
		bookName := htmlquery.SelectAttr(bookNameMeta, "content")
		fmt.Println("书名 = ", bookName)

		//获取书作者
		AuthorMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:novel:author']")
		author := htmlquery.SelectAttr(AuthorMeta, "content")
		fmt.Println("作者 = ", author)

		//获取书的描述信息
		DescriptionMeta, _ := htmlquery.FindOne(doc, "//meta[@property='og:description']")
		description := htmlquery.SelectAttr(DescriptionMeta, "content")
		fmt.Println("简介 = ", description)

		//获取书章节列表
		ddNode, _ := htmlquery.Find(doc, "//div[@id='list']//dl//dd")
		for i := 0; i < len(ddNode); i++ {
			var tmp Chapter
			aNode, _ := htmlquery.Find(ddNode[i], "//a")
			tmp.Link = "https://www.999xs.com" + htmlquery.SelectAttr(aNode[0], "href")
			tmp.Title = htmlquery.InnerText(aNode[0])
			chapters = append(chapters, tmp)
		}

		//导入信息
		bi = BookInfo{
			EBHost:      this.URL,
			EBookID:     bookid,
			Name:        bookName,
			Author:      author,
			Description: description,
			Chapters:    chapters,
		}
	}
	return bi
}
func (this Ebook999XS) DownloadChapters(Bi BookInfo, proxy string) BookInfo {
	result := Bi //先进行赋值，把数据
	var chapters []Chapter
	result.Chapters = chapters //把原来的数据清空
	bis := Bi.Split()

	for index := 0; index < len(bis); index++ {
		this.Lock.Lock()
		bookinfo := bis[index]
		rec := this.downloadChapters(bookinfo, "")
		chapters = append(chapters, rec.Chapters...)
		//fmt.Printf("Get into this.Lock.Unlock() time: %d\n", index+1)
		this.Lock.Unlock()
	}
	result.Chapters = chapters

	return result
}

//根据每个章节的 URL连接，下载每章对应的内容Content当中
func (this Ebook999XS) downloadChapters(Bi BookInfo, proxy string) BookInfo {
	chapters := Bi.Chapters

	NumChapter := len(chapters)
	tmpChapter := make(chan Chapter, NumChapter)
	ResultCh := make(chan chan Chapter, NumChapter)
	wg := sync.WaitGroup{}
	var c []Chapter
	var bar *progressbar.ProgressBar
	go AsycChapter(ResultCh, tmpChapter)
	for index := 0; index < NumChapter; index++ {
		tmp := ProxyChapter{
			Proxy: proxy,
			C:     chapters[index],
		}
		this.DownloaderChapter(ResultCh, tmp, &wg)
	}

	wg.Wait()

	//下载章节的时候显示进度条
	bar = progressbar.NewOptions(
		NumChapter,
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowIts(),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "#", SaucerPadding: "-", BarStart: ">", BarEnd: "<"}),
	)

	for index := 0; index <= NumChapter; {
		select {
		case tmp := <-tmpChapter:
			//fmt.Printf("tmp.Title = %s\n", tmp.Title)
			//fmt.Printf("tmp.Content= %s\n", tmp.Content)
			c = append(c, tmp)
			index++
			if index == NumChapter {
				goto ForEnd
			}
		}
		bar.Add(1)

	}
ForEnd:

	result := BookInfo{
		EBHost:      Bi.EBHost,
		EBookID:     Bi.EBookID,
		Name:        Bi.Name,
		Author:      Bi.Author,
		Description: Bi.Description,
		Volumes:     Bi.Volumes,       //小说分卷信息在 GetBookInfo()的时候已经下载完成
		HasVolume:   Bi.VolumeState(), //小说分卷信息在 GetBookInfo()的时候已经定义
		Chapters:    c,
	}

	return result
}

//DownloaderChapter 下载小说
func (this Ebook999XS) DownloaderChapter(ResultChan chan chan Chapter, pc ProxyChapter, wg *sync.WaitGroup) {
	c := make(chan Chapter)
	ResultChan <- c
	wg.Add(1)
	go func(pc ProxyChapter) {
		pollURL := pc.C.Link
		proxy := pc.Proxy
		var result Chapter

		if proxy != "" {
			doc, _ := htmlquery.LoadURLWithProxy(pollURL, proxy)
			contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
			contentText := htmlquery.OutputHTML(contentNode, false)

			//替换两个 html换行
			tmp := strings.Replace(contentText, "<br/><br/>", "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "<br/>", "\r\n", -1)

			//把 readx(); 替换成 ""
			tmp = strings.Replace(tmp, "999小说更新最快 电脑端:https://www.999xs.com/", "", -1)
			tmp = strings.Replace(tmp, "ωωω.九九九xs.com", "", -1)
			tmp = strings.Replace(tmp, "999小说首发 https://www.999xs.com https://m.999xs.com", "", -1)
			tmp = strings.Replace(tmp, "手机\\端 一秒記住『www.999xs.com』為您提\\供精彩小說\\閱讀", "", -1)

			//tmp = tmp + "\r\n"
			//返回数据，填写Content内容
			result = Chapter{
				Title:   pc.C.Title,
				Link:    pc.C.Link,
				Content: tmp,
			}
		} else {
			doc, _ := htmlquery.LoadURL(pollURL)
			contentNode, _ := htmlquery.FindOne(doc, "//div[@id='content']")
			contentText := htmlquery.OutputHTML(contentNode, false)

			//替换两个 html换行
			tmp := strings.Replace(contentText, "<br/><br/>", "\r\n", -1)
			//替换一个 html换行
			tmp = strings.Replace(tmp, "<br/>", "\r\n", -1)

			//把 readx(); 替换成 ""
			tmp = strings.Replace(tmp, "999小说更新最快 电脑端:https://www.999xs.com/", "", -1)
			tmp = strings.Replace(tmp, "ωωω.九九九xs.com", "", -1)
			tmp = strings.Replace(tmp, "999小说首发 https://www.999xs.com https://m.999xs.com", "", -1)
			tmp = strings.Replace(tmp, "手机\\端 一秒記住『www.999xs.com』為您提\\供精彩小說\\閱讀", "", -1)

			//tmp = tmp + "\r\n"
			//返回数据，填写Content内容
			result = Chapter{
				Title:   pc.C.Title,
				Link:    pc.C.Link,
				Content: tmp,
			}
		}
		c <- result
		wg.Done()
	}(pc)
}
