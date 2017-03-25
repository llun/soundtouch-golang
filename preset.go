package soundtouch

type Preset struct {
  ID      int         `xml:"id,attr"`
  Content ContentItem `xml:"ContentItem"`
}
