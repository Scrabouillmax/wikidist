<template>
  <div id="graph" ref="svg">
    <svg v-show="!loading" height="100%" width="100%">
      <g :transform="centerTranslate">
        <line
          v-for="(link, index) of links"
          :key="index"
          stroke="#75b1bf"
          stroke-width="2"
          :x1="nodes[link.source.index].x * 20"
          :y1="nodes[link.source.index].y * 20"
          :x2="nodes[link.target.index].x * 20"
          :y2="nodes[link.target.index].y * 20"
        />
        <g
          v-for="node of nodes"
          :key="node.article.uid"
          :transform="`translate(${node.x * 20} ${node.y * 20})`"
        >
          <circle r="10" fill="#f47370"></circle>
          <a
            :href="`https://simple.wikipedia.org/wiki/${node.article.title.replace(' ', '_')}`"
            target="_blank"
          >
            <text x="15" font-size="0.8em">{{ node.article.title }}</text>
          </a>
        </g>
      </g>
    </svg>
    <img v-show="loading" id="loader" src="@/assets/loader.gif" />
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator';
import { forceManyBody, forceSimulation, Simulation, forceLink, forceCenter } from 'd3-force';
import { ArticleNode, ArticleLink, Article } from '@/types/article';

@Component
export default class Graph extends Vue {
  @Prop() loading!: boolean;
  @Prop() links!: ArticleLink[];
  @Prop() articles!: Article[];

  private height: number = 0;
  private width: number = 0;

  private nodes: ArticleNode[] = [];

  mounted() {
    this.height = (this.$refs.svg as Element).clientHeight;
    this.width = (this.$refs.svg as Element).clientWidth;
  }

  @Watch('articles')
  public onArticleChange() {
    this.updateNodes();
  }

  updateNodes() {
    const newNodes: ArticleNode[] = [];
    for (const node of this.nodes) {
      if (this.articles.find(a => a.uid === node.article.uid)) {
        newNodes.push(node);
      }
    }
    this.nodes = newNodes;

    for (const article of this.articles) {
      if (!this.nodes.find(n => n.article.uid === article.uid)) {
        this.nodes.push({
          article,
          x: 0,
          y: 0,
        });
      }
    }
    this.simulation
      .nodes(this.nodes)
      .force('charge', forceManyBody().strength(-3000))
      .force(
        'links',
        forceLink(this.links)
          .id(node => (node as ArticleNode).article.uid)
          .distance(100)
          .strength(100)
      )
      .force('center', forceCenter())
      .restart();
  }

  public get centerTranslate() {
    return `translate(${this.width / 2} ${this.height / 2})`;
  }

  private simulation: Simulation<ArticleNode, ArticleLink> = forceSimulation(this.nodes)
    .force('charge', forceManyBody().strength(-3000))
    .force(
      'links',
      forceLink(this.links)
        .id(node => (node as ArticleNode).article.uid)
        .distance(200)
        .strength(200)
    )
    .force('center', forceCenter());
}
</script>
<style>
#loader {
  height: 300px;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
}

#graph {
  height: calc(100% - 50px);
  display: flex;
}
</style>
