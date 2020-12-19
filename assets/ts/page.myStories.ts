const myStories = document.getElementById("myStories");

const myStoriesDraftTab = myStories?.querySelector(
  ".myStories-tab__drafts"
) as HTMLElement | null;
const myStoriesPublishedTab = myStories?.querySelector(
  ".myStories-tab__published"
) as HTMLElement | null;

const myDraftStoriesWrapper = myStories?.querySelector(
  ".myStories__wrapper.myStories-drafts"
) as HTMLElement | null;
const myPublishedStoriesWrapper = myStories?.querySelector(
  ".myStories__wrapper.myStories-published"
) as HTMLElement | null;

const init = () => {
  myStoriesDraftTab?.addEventListener("click", (e) => {
    const currentTarget = e.currentTarget as HTMLElement | null;
    if (
      !currentTarget ||
      !myStoriesPublishedTab ||
      !myDraftStoriesWrapper ||
      !myPublishedStoriesWrapper
    ) {
      return;
    }
    currentTarget.style.borderBottomColor = "transparent";
    myStoriesPublishedTab.style.borderBottomColor = "rgba(0, 0, 0, 0.9)";
    myDraftStoriesWrapper.className = myDraftStoriesWrapper.className.replace(
      "none",
      "block"
    );
    myPublishedStoriesWrapper.className = myPublishedStoriesWrapper.className.replace(
      "block",
      "none"
    );
  });
  myStoriesPublishedTab?.addEventListener("click", (e) => {
    const currentTarget = e.currentTarget as HTMLElement | null;
    if (
      !currentTarget ||
      !myStoriesDraftTab ||
      !myDraftStoriesWrapper ||
      !myPublishedStoriesWrapper
    ) {
      return;
    }
    currentTarget.style.borderBottomColor = "transparent";
    myStoriesDraftTab.style.borderBottomColor = "rgba(0, 0, 0, 0.9)";
    myDraftStoriesWrapper.className = myDraftStoriesWrapper.className.replace(
      "block",
      "none"
    );
    myPublishedStoriesWrapper.className = myPublishedStoriesWrapper.className.replace(
      "none",
      "block"
    );
  });
};

if (document.location.pathname.includes("/me/stories")) {
  init();
}
