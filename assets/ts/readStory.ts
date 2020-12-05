import EditorJS from "@editorjs/editorjs";
import Header from "@editorjs/header";
import CodeTool from "@editorjs/code";
import ImageTool from "@editorjs/image";
import { BASE_URL } from "./constants";
import Axios from "axios";

const editorReadOnlyHeader = document.getElementById("editor-readOnly__header");

const initEditorReadOnly = async (storyId: string) => {
  const { data: blocks } = await Axios.get(BASE_URL + `/blocks/${storyId}`);
  const header = blocks.shift();
  const headerEditor = new EditorJS({
    //   readOnly: true, // This occurs error on 2.19.0 version. It's on github issue https://github.com/codex-team/editor.js/issues/1400
    holder: "editor-readOnly__header",
    tools: {
      header: {
        class: Header,
        config: {
          levels: [2, 4, 6],
        },
      },
      code: CodeTool,
      image: {
        class: ImageTool,
      },
    },
    data: { blocks: [header] },
  });
  const bodyEditor = new EditorJS({
    //   readOnly: true, // This occurs error on 2.19.0 version. It's on github issue https://github.com/codex-team/editor.js/issues/1400
    holder: "editor-readOnly__body",
    tools: {
      header: {
        class: Header,
        config: {
          levels: [2, 4, 6],
        },
      },
      code: CodeTool,
      image: {
        class: ImageTool,
      },
    },
    data: { blocks },
  });
  headerEditor.isReady.then(async () => {
    await headerEditor.readOnly.toggle(true);
    const x = editorReadOnlyHeader?.querySelector(".codex-editor__redactor") as
      | HTMLElement
      | null
      | undefined;
    console.log(x);
    x!.style.paddingBottom = "1rem";
  });
  bodyEditor.isReady.then(async () => {
    await bodyEditor.readOnly.toggle(true);
  });
};

const readStory = async () => {
  if (BASE_URL) {
    const params = document.location.pathname.split(BASE_URL)[0].split("/");
    if (params[1] === "read") {
      const storyId = params[2];
      await initEditorReadOnly(storyId);
    }
  }
};

readStory();
