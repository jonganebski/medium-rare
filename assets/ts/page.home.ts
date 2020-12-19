import Axios from "axios";
import { BASE_URL, MONTHS } from "./constants";

const homeMainList = document.getElementById(
  "home__mainList"
) as HTMLElement | null;

const requestMoreStoryCards = async () => {
  if (
    window.innerHeight + window.pageYOffset >
    document.body.offsetHeight * 0.9
  ) {
    window.removeEventListener("scroll", requestMoreStoryCards);
    const lastStoryTimestamp = homeMainList?.lastElementChild?.id;
    try {
      const { status, data } = await Axios.get(
        `/api/recent-stories/${lastStoryTimestamp}`
      );
      if (status < 300 && data?.length !== 0) {
        data.forEach((story: any) => {
          const storyCard = drawStoryCard(story);
          homeMainList?.append(storyCard);
        });
        window.addEventListener("scroll", requestMoreStoryCards);
      }
    } catch {}
  }
};

const formatPostDate = (timeInSec: number): string => {
  const now = Math.round(Date.now() / 1000);
  const lapse = now - timeInSec;
  const oneDayInSec = 24 * 60 * 60;
  if (lapse < oneDayInSec) {
    return "today";
  }
  if (lapse < 2 * oneDayInSec) {
    return "yesterday";
  }
  if (lapse < 3 * oneDayInSec) {
    return "2 days ago";
  }
  const dateObj = new Date(timeInSec);
  const month = MONTHS[dateObj.getMonth()];
  const date = dateObj.getDate();
  const year = dateObj.getFullYear();
  return `${month} ${date}, ${year}`;
  // return time.Unix(createdAt, 0).Format("January 2, 2006")
};

const grindBody = (body: string, length: number): string => {
  return body.slice(0, length);
};

const drawStoryCard = (story: any) => {
  const coverImgEl = document.createElement("img");
  const coverFrameEl = document.createElement("a");
  const readTimeEl = document.createElement("span");
  const dotEl = document.createElement("div");
  const createdAtEl = document.createElement("span");
  const infoEl = document.createElement("div");
  const authorNameEl = document.createElement("a");
  const containerEl = document.createElement("div");
  const bodyEl = document.createElement("p");
  const headerEl = document.createElement("h3");
  const referenceEl = document.createElement("div");
  const linkToReadEl = document.createElement("a");
  const contentEl = document.createElement("div");
  const liEl = document.createElement("li");
  coverImgEl.className = "storyCard__cover-img";
  coverImgEl.src = story.coverImgUrl;
  coverFrameEl.className = "storyCard__cover-frame";
  coverFrameEl.href = `/read-story/${story.storyId}`;
  readTimeEl.innerText = story.readTime;
  dotEl.className = "_devider-dot";
  createdAtEl.innerText = formatPostDate(story.createdAt);
  infoEl.className = "storyCard__info _flex-cs";
  authorNameEl.className = "storyCard__author-name _block";
  authorNameEl.href = `/user-home/${story.authorId}`;
  authorNameEl.innerText = story.authorUsername;
  bodyEl.className = "storyCard__body";
  bodyEl.innerText = grindBody(story.body, 120);
  headerEl.className = "storyCard__header";
  headerEl.innerText = story.header;
  referenceEl.className = "storyCard__reference";
  referenceEl.innerText = "NOT BASED ON YOUR READING HISTORY";
  linkToReadEl.href = `/read-story/${story.storyId}`;
  contentEl.className = "storyCard__content _flex-c-sb";
  liEl.className = "storyCard _flex";
  liEl.id = story.createdAt;
  coverFrameEl.append(coverImgEl);
  infoEl.append(createdAtEl);
  infoEl.append(dotEl);
  infoEl.append(readTimeEl);
  containerEl.append(authorNameEl);
  containerEl.append(infoEl);
  linkToReadEl.append(referenceEl);
  linkToReadEl.append(headerEl);
  linkToReadEl.append(bodyEl);
  contentEl.append(linkToReadEl);
  contentEl.append(containerEl);
  liEl.append(contentEl);
  liEl.append(coverFrameEl);
  return liEl;
};

const init = () => {
  window.addEventListener("scroll", requestMoreStoryCards);
};

if (BASE_URL && document.location.pathname.split(BASE_URL)[0] === "/") {
  init();
}
