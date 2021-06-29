package ics

import (
	"github.com/tongruirenye/OrgICSX5/server/org"
)

type IcsWriter struct {
	extendingWriter org.Writer
}

func NewIcsWriter() *IcsWriter {
	return &IcsWriter{}
}

func (w *IcsWriter) Before(d *org.Document) {}

func (w *IcsWriter) After(d *org.Document) {}

func (w *IcsWriter) String() string {
	return ""
}

func (w *IcsWriter) WriterWithExtensions() org.Writer {
	return w
}

func (w *IcsWriter) WriteNodesAsString(nodes ...org.Node) string {
	return ""
}

func (w *IcsWriter) WriteKeyword(k org.Keyword) {

}
func (w *IcsWriter) WriteInclude(i org.Include) {

}
func (w *IcsWriter) WriteComment(c org.Comment) {

}
func (w *IcsWriter) WriteNodeWithMeta(n org.NodeWithMeta) {

}
func (w *IcsWriter) WriteNodeWithName(n org.NodeWithName) {

}
func (w *IcsWriter) WriteHeadline(h org.Headline) {

}
func (w *IcsWriter) WriteBlock(b org.Block) {

}
func (w *IcsWriter) WriteResult(r org.Result) {

}
func (w *IcsWriter) WriteInlineBlock(i org.InlineBlock) {

}
func (w *IcsWriter) WriteExample(e org.Example) {

}
func (w *IcsWriter) WriteDrawer(d org.Drawer) {

}
func (w *IcsWriter) WritePropertyDrawer(p org.PropertyDrawer) {

}
func (w *IcsWriter) WriteList(l org.List) {

}
func (w *IcsWriter) WriteListItem(l org.ListItem) {

}
func (w *IcsWriter) WriteDescriptiveListItem(d org.DescriptiveListItem) {

}
func (w *IcsWriter) WriteTable(t org.Table) {

}
func (w *IcsWriter) WriteHorizontalRule(h org.HorizontalRule) {

}
func (w *IcsWriter) WriteParagraph(p org.Paragraph) {

}
func (w *IcsWriter) WriteText(t org.Text) {

}
func (w *IcsWriter) WriteEmphasis(e org.Emphasis) {

}
func (w *IcsWriter) WriteLatexFragment(l org.LatexFragment) {

}
func (w *IcsWriter) WriteStatisticToken(s org.StatisticToken) {

}
func (w *IcsWriter) WriteExplicitLineBreak(e org.ExplicitLineBreak) {

}
func (w *IcsWriter) WriteLineBreak(l org.LineBreak) {

}
func (w *IcsWriter) WriteRegularLink(r org.RegularLink) {

}
func (w *IcsWriter) WriteMacro(m org.Macro) {

}
func (w *IcsWriter) WriteTimestamp(t org.Timestamp) {

}
func (w *IcsWriter) WriteFootnoteLink(f org.FootnoteLink) {

}
func (w *IcsWriter) WriteFootnoteDefinition(f org.FootnoteDefinition) {

}

func (w *IcsWriter) WriteTimeProperty(f org.TimeProperty) {

}
