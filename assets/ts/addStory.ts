import EditorJS from "@editorjs/editorjs";
import Header from "@editorjs/header";
import CodeTool from "@editorjs/code";
import ImageTool from "@editorjs/image";

// const textEditor = document.getElementById("editor__container");
// const textEditorTextBoxes = textEditor?.querySelectorAll(".textBox");

const editor = new EditorJS({
  holder: "editor__container",
  tools: {
    header: {
      class: Header,
      inlineToolbar: true,
    },
    code: CodeTool,
    image: {
      class: ImageTool,
      config: {
        endpoints: {
          byFile: "http://localhost:4000/upload/photo/byfile",
        },
      },
    },
  },
});

const addStory = () => {
  // document.body.addEventListener("click", () => {
  //   editor.save().then((savedData) => console.log(savedData));
  // });
};

addStory();
