package main

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
   <head>
      <meta charset="utf-8" />
      <meta name="viewport" content="width=device-width,initial-scale=1" />
      <title>TweetMe Card</title>
      <style>
         :root{
         --bg: #f9fafb;
         --card-bg: #ffffff;
         --text: #000000;
         --muted: #6b7280;
         --border: #e5e7eb;
         --border-strong: #e5e7eb;
         --link: #60a5fa;
         }
         .dark{ --bg: #000000; --card-bg: #1f2937; --text: #ffffff; --muted: #9ca3af; --border: #4b5563; --border-strong: #1f2937; --link: #ffffff; }
         * { box-sizing: border-box; }
         html, body {
         font-size: 40px;
         height: 100%; margin: 0; 
         }
         body{ margin:0; font-family: ui-sans-serif, system-ui, -apple-system, Segoe UI, Roboto, Noto Sans,
         Ubuntu, Cantarell, Helvetica Neue, Arial, "Apple Color Emoji",
         "Segoe UI Emoji", "Segoe UI Symbol"; background: var(--bg); color: var(--text); line-height: 1.5;
		 background: linear-gradient(90deg, hsla(212, 35%, 58%, 1) 0%, hsla(218, 32%, 80%, 1) 100%);
         }
         .page{ width: 1600px; height: 900px; display: flex; align-items: center; justify-content: center; padding: 2.5rem; margin: 0 auto; }
         .tweet{ background: var(--card-bg); border: 2px solid var(--border); border-radius: 0.75rem; padding: 1rem 1rem 0.5rem; width: 100%; max-width: 1400px; box-shadow: 0 10px 30px rgba(2,6,23,0.08); }
         .row{ display:flex; justify-content: space-between; align-items: center; gap: 0.75rem; }
         .author{ display:flex; align-items:center; gap: 0.75rem; }
         .avatar{ width: 3rem; height: 3rem; border-radius: 9999px; background: var(--border); object-fit: cover; margin-right: 0.5rem; flex: 0 0 auto; }
         .name{ display:block; font-weight: 700; font-size: 0.95rem; line-height: 1.2; color: var(--text); margin: 0; }
         .handle{ display:block; font-weight: 400; font-size: 0.9rem; line-height: 1.2; color: var(--muted); margin: 2px 0 0; }
         .content{ margin-top: 0.75rem; font-size: 1.25rem; line-height: 1.35; color: var(--text); }
         .meta{ color: var(--muted); font-size: 1rem; padding: 0.25rem 0; margin: 0 0.125rem; }
         .divider{ border: 1px solid var(--border); border-bottom: 0; margin: 0.25rem 0; }
         .actions{ display:flex; align-items: center; margin-top: 0.25rem; color: var(--muted); flex-wrap: wrap; }
         .act{ display:flex; align-items:center; margin-right: 1.5rem; }
         .act svg{ width: 1rem; height: 1rem; fill: none; margin-right: 0.5rem; flex: 0 0 auto; }
      </style>
   </head>
   <body class="{{ .BodyClass }}">
      <div class="page">
         <article class="tweet" role="article" aria-label="Tweet">
            <div class="row">
               <div class="author">
                  <img class="avatar" alt="{{ .User }} avatar" src="{{ .Avatar }}" />
                  <div>
                     <p class="name">
                        {{ .User }}
                        <svg xmlns="http://www.w3.org/2000/svg" width="1.2em" height="1.2em" viewBox="0 0 24 24" fill="#1D9BF0" style="vertical-align: middle; margin-left: 0.25em;">
                           <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
                           <path d="M12.01 2.011a3.2 3.2 0 0 1 2.113 .797l.154 .145l.698 .698a1.2 1.2 0 0 0 .71 .341l.135 .008h1a3.2 3.2 0 0 1 3.195 3.018l.005 .182v1c0 .27 .092 .533 .258 .743l.09 .1l.697 .698a3.2 3.2 0 0 1 .147 4.382l-.145 .154l-.698 .698a1.2 1.2 0 0 0 -.341 .71l-.008 .135v1a3.2 3.2 0 0 1 -3.018 3.195l-.182 .005h-1a1.2 1.2 0 0 0 -.743 .258l-.1 .09l-.698 .697a3.2 3.2 0 0 1 -4.382 .147l-.154 -.145l-.698 -.698a1.2 1.2 0 0 0 -.71 -.341l-.135 -.008h-1a3.2 3.2 0 0 1 -3.195 -3.018l-.005 -.182v-1a1.2 1.2 0 0 0 -.258 -.743l-.09 -.1l-.697 -.698a3.2 3.2 0 0 1 -.147 -4.382l.145 -.154l.698 -.698a1.2 1.2 0 0 0 .341 -.71l.008 -.135v-1l.005 -.182a3.2 3.2 0 0 1 3.013 -3.013l.182 -.005h1a1.2 1.2 0 0 0 .743 -.258l.1 -.09l.698 -.697a3.2 3.2 0 0 1 2.269 -.944zm3.697 7.282a1 1 0 0 0 -1.414 0l-3.293 3.292l-1.293 -1.292l-.094 -.083a1 1 0 0 0 -1.32 1.497l2 2l.094 .083a1 1 0 0 0 1.32 -.083l4 -4l.083 -.094a1 1 0 0 0 -.083 -1.32z" />
                        </svg>
                     </p>
                     <span class="handle">@{{ .Handle }}</span>
                  </div>
               </div>
            </div>
            <p class="content">{{ .Text }}</p>
            <div class="divider" aria-hidden="true"></div>
            <div class="actions" aria-label="Tweet actions">
               <div class="act" aria-label="Quotes">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2.992 16.342a2 2 0 0 1 .094 1.167l-1.065 3.29a1 1 0 0 0 1.236 1.168l3.413-.998a2 2 0 0 1 1.099.092 10 10 0 1 0-4.777-4.719"/></svg>
                  <span>{{ .Quotes }}</span>
               </div>
               <div class="act" aria-label="Retweet">
                  <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-repeat"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M4 12v-3a3 3 0 0 1 3 -3h13m-3 -3l3 3l-3 3" /><path d="M20 12v3a3 3 0 0 1 -3 3h-13m3 3l-3 -3l3 -3" /></svg>
                  <span>{{ .Retweets }}</span>
               </div>
               <div class="act" aria-label="Likes">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 9.5a5.5 5.5 0 0 1 9.591-3.676.56.56 0 0 0 .818 0A5.49 5.49 0 0 1 22 9.5c0 2.29-1.5 4-3 5.5l-5.492 5.313a2 2 0 0 1-3 .019L5 15c-1.5-1.5-3-3.2-3-5.5"/></svg>
                  <span>{{ .Likes }}</span>
               </div>
               <p class="meta">{{ .Timestamp }}</p>
            </div>
         </article>
      </div>
   </body>
</html>
`
