use std::borrow::Cow;

use html5ever::parse_document;
use html5ever::tendril::*;
use html5ever::tree_builder::{AppendText, ElementFlags, NodeOrText, QuirksMode, TreeSink};
use html5ever::{Attribute, ExpandedName, QualName};
use std::collections::HashMap;

#[derive(Default)]
struct Sink {
    title: Option<String>,
    title_node: Option<usize>,
    next_id: usize,
    names: HashMap<usize, QualName>,
}

impl Sink {
    fn get_id(&mut self) -> usize {
        let id = self.next_id;
        self.next_id += 1;
        id
    }
}

impl TreeSink for Sink {
    type Handle = usize;
    type Output = Option<String>;

    fn finish(self) -> Self::Output {
        self.title
    }

    fn parse_error(&mut self, _msg: Cow<'static, str>) {}

    fn get_document(&mut self) -> Self::Handle {
        0
    }

    fn get_template_contents(&mut self, target: &Self::Handle) -> Self::Handle {
        if let Some(expanded_name!(html "template")) = self.names.get(target).map(|n| n.expanded())
        {
            target + 1
        } else {
            panic!("not a template element")
        }
    }

    fn set_quirks_mode(&mut self, _mode: QuirksMode) {}

    fn same_node(&self, x: &Self::Handle, y: &Self::Handle) -> bool {
        x == y
    }

    fn elem_name<'a>(&'a self, target: &'a Self::Handle) -> ExpandedName<'a> {
        self.names.get(target).expect("not an element").expanded()
    }

    fn create_element(
        &mut self,
        name: QualName,
        _attrs: Vec<Attribute>,
        _flags: ElementFlags,
    ) -> Self::Handle {
        let id = self.get_id();
        if name.expanded() == expanded_name!(html "title") {
            self.title_node = Some(id);
        }
        self.names.insert(id, name);
        id
    }

    fn create_comment(&mut self, _text: StrTendril) -> Self::Handle {
        self.get_id()
    }

    #[allow(unused_variables)]
    fn create_pi(&mut self, target: StrTendril, data: StrTendril) -> Self::Handle {
        unimplemented!()
    }

    fn append(&mut self, parent: &Self::Handle, child: NodeOrText<Self::Handle>) {
        if Some(*parent) == self.title_node {
            if let AppendText(t) = child {
                self.title = Some(t.to_string());
            }
        }
    }

    fn append_before_sibling(
        &mut self,
        _sibling: &Self::Handle,
        _new_node: NodeOrText<Self::Handle>,
    ) {
    }

    fn append_based_on_parent_node(
        &mut self,
        element: &Self::Handle,
        _prev_element: &Self::Handle,
        child: NodeOrText<Self::Handle>,
    ) {
        self.append_before_sibling(element, child);
    }

    fn append_doctype_to_document(
        &mut self,
        _name: StrTendril,
        _public_id: StrTendril,
        _system_id: StrTendril,
    ) {
    }

    fn add_attrs_if_missing(&mut self, _target: &Self::Handle, _attrs: Vec<Attribute>) {}

    fn remove_from_parent(&mut self, _target: &Self::Handle) {}

    fn reparent_children(&mut self, _node: &Self::Handle, _new_parent: &Self::Handle) {}

    fn mark_script_already_started(&mut self, _node: &Self::Handle) {}

    fn set_current_line(&mut self, _line_number: u64) {}

    fn pop(&mut self, _node: &Self::Handle) {}
}

pub fn parse<R: std::io::Read>(s: &mut R) -> Option<String> {
    let sink = Sink {
        next_id: 1,
        title: None,
        title_node: None,
        names: HashMap::new(),
    };
    parse_document(sink, Default::default())
        .from_utf8()
        .read_from(s)
        .unwrap_or_else(|_| None)
}

#[cfg(test)]
mod test {
    use super::*;
    #[test]
    fn html() {
        let title = r#"
            <!DOCTYPE html>
            <html lang="en">
            <head>
              <meta charset="utf-8">
              <title>Hello World!</title>
            </head>
            <body>
              <div id="root"></div>
            </body>
        "#;
        assert_eq!(parse(&mut std::io::Cursor::new(title)), Some("Hello World!".to_owned()));
    }
    #[test]
    fn xhtml() {
        let title = r#"
            <?xml version="1.0" encoding="UTF-8"?>

            <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
                "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">

            <html xmlns="http://www.w3.org/TR/xhtml1" xml:lang="en" lang="en">
           <head>
              <title>Hello World!</title>
           </head>
           <body>
              <div id="root"></div>
           </body>
    </html>
        "#;
        assert_eq!(parse(&mut std::io::Cursor::new(title)), Some("Hello World!".to_owned()));
    }

    #[test]
    fn broken() {
        let title = r#"
            DOCTYPE html>
            <html lang "en">
            <head>
              <meta charset="utf-8">
              <title>Hello World!</title>
            </head>
            <body>
              <div id="root">
            </bo>
        "#;
        assert_eq!(parse(&mut std::io::Cursor::new(title)), Some("Hello World!".to_owned()));
    }


    #[test]
    fn notitle() {
        let title = r#"
            DOCTYPE html>
            <html lang "en">
            <head>
              <meta charset="utf-8">
            </head>
            <body>
              <div id="root">
            </body>
        "#;
        assert_eq!(parse(&mut std::io::Cursor::new(title)), None);
    }
}
