<script setup>
import {computed, h, onBeforeMount, ref} from 'vue'
import {
  GetAIResponseResult,
  GetConfig, GetPromptTemplates,
  GetTelegraphList,
  GlobalStockIndexes,
  SaveAIResponseResult, SaveAsMarkdown, ShareAnalysis,
  SummaryStockNews
} from "../../wailsjs/go/main/App";
import {EventsOn} from "../../wailsjs/runtime";
import NewsList from "./newsList.vue";
import KLineChart from "./KLineChart.vue";
import {Add, ChatboxOutline, PulseOutline,} from "@vicons/ionicons5";
import {NAvatar, NButton, NFlex, NText, useMessage, useNotification} from "naive-ui";
import {ExportPDF} from "@vavt/v3-extension";
import {MdEditor, MdPreview} from "md-editor-v3";
const icon = ref('https://raw.githubusercontent.com/ArvinLovegood/go-stock/master/build/appicon.png');

const message = useMessage()
const notify = useNotification()
const panelHeight = ref(window.innerHeight - 240)

const telegraphList = ref([])
const sinaNewsList = ref([])

const common = ref([])
const america = ref([])
const europe = ref([])
const asia = ref([])
const other = ref([])
const globalStockIndexes = ref(null)
const summaryModal= ref(false)
const summaryBTN= ref(true)
const darkTheme= ref(false)
const theme=computed(() => {
  return darkTheme ? 'dark' : 'light'
})
const aiSummary=ref(``)
const aiSummaryTime=ref("")
const modelName=ref("")
const chatId=ref("")
const question=ref(``)
const sysPromptId=ref(0)
const loading=ref(true)
const sysPromptOptions=ref([])
const userPromptOptions=ref([])
const promptTemplates=ref([])

function getIndex() {
  GlobalStockIndexes().then((res) => {
    globalStockIndexes.value = res
    common.value = res["common"]
    america.value = res["america"]
    europe.value = res["europe"]
    asia.value = res["asia"]
    other.value = res["other"]
  })
}

onBeforeMount(() => {
  GetConfig().then(result => {
      summaryBTN.value= result.openAiEnable
      darkTheme.value = result.darkTheme
  })
  GetPromptTemplates("","").then(res=>{
    promptTemplates.value=res
    sysPromptOptions.value=promptTemplates.value.filter(item => item.type === '模型系统Prompt')
    userPromptOptions.value=promptTemplates.value.filter(item => item.type === '模型用户Prompt')
  })

  GetTelegraphList("财联社电报").then((res) => {
    telegraphList.value = res
  })
  GetTelegraphList("新浪财经").then((res) => {
    sinaNewsList.value = res
  })
  getIndex();

  setInterval(() => {
    getIndex()
  }, 3000)
})



EventsOn("newTelegraph", (data) => {
  for (let i = 0; i < data.length; i++) {
    telegraphList.value.pop()
  }
  telegraphList.value.unshift(...data)
})
EventsOn("newSinaNews", (data) => {
  for (let i = 0; i < data.length; i++) {
    sinaNewsList.value.pop()
  }
  sinaNewsList.value.unshift(...data)
})

//获取页面高度
window.onresize = () => {
  panelHeight.value = window.innerHeight - 240
}

function getAreaName(code){
  switch (code) {
    case "america":
      return "美国"
    case "europe":
      return "欧洲"
    case "asia":
      return "亚洲"
    case "common":
      return "常用"
    case "other":
      return "其他"
  }
}
function reAiSummary(){
  aiSummary.value=""
  summaryModal.value = true
  loading.value = true
  SummaryStockNews(question.value,sysPromptId.value)
}
function getAiSummary(){
  summaryModal.value = true
  loading.value = true
  GetAIResponseResult("市场资讯").then(result => {
    if(result.content){
      aiSummary.value=result.content
      question.value=result.question
      loading.value = false

      const date = new Date(result.CreatedAt);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      aiSummaryTime.value=`${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
      modelName.value=result.modelName
    }else{
      aiSummaryTime.value=""
      aiSummary.value=""
      modelName.value=""
      SummaryStockNews(question.value,sysPromptId.value)
    }
  })
}

function updateTab(name) {
  summaryBTN.value = (name === "市场快讯");
}

EventsOn("summaryStockNews",async (msg) => {
  loading.value = false
  ////console.log(msg)
  if (msg === "DONE") {
    SaveAIResponseResult("市场资讯","市场资讯", aiSummary.value, chatId.value,question.value)
    message.info("AI分析完成！")
    message.destroyAll()

  } else {
    if(msg.chatId){
      chatId.value = msg.chatId
    }
    if(msg.question){
      question.value = msg.question
    }
    if(msg.content){
      aiSummary.value =aiSummary.value + msg.content
    }
    if(msg.extraContent){
      aiSummary.value = aiSummary.value + msg.extraContent
    }
    if(msg.model){
      modelName.value=msg.model
    }
    if(msg.time){
      aiSummaryTime.value = msg.time
    }
  }
})

async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(aiSummary.value);
    message.success('分析结果已复制到剪切板');
  } catch (err) {
    message.error('复制失败: ' + err);
  }
}
function saveAsMarkdown(){
  SaveAsMarkdown('市场资讯','市场资讯').then(result => {
    message.success(result)
  })
}
function share(){
  ShareAnalysis('市场资讯','市场资讯').then(msg => {
    //message.info(msg)
    notify.info({
      avatar: () =>
          h(NAvatar, {
            size: 'small',
            round: false,
            src: icon.value
          }),
      title: '分享到社区',
      duration:1000*30,
      content: () => {
        return h('div', {
          style: {
            'text-align': 'left',
            'font-size': '14px',
          }
        }, { default: () => msg })
      },
    })
  })
}
</script>

<template>
  <n-card>
    <n-tabs type="line" animated @update-value="updateTab">
      <n-tab-pane name="市场快讯" tab="市场快讯"  >
        <n-grid :cols="2" :y-gap="0">
          <n-gi>
            <news-list :newsList="telegraphList" :header-title="'财联社电报'"></news-list>
          </n-gi>
          <n-gi>
            <news-list :newsList="sinaNewsList" :header-title="'新浪财经'"></news-list>
          </n-gi>
        </n-grid>
      </n-tab-pane>
      <n-tab-pane name="全球股指" tab="全球股指">
        <n-tabs type="segment" animated>
          <n-tab-pane name="全球指数" tab="全球指数">
            <n-grid :cols="5" :y-gap="0">
              <n-gi v-for="(val, key) in globalStockIndexes" :key="key">
                <n-list bordered>
                  <template #header>
                    {{ getAreaName(key) }}
                  </template>
                  <n-list-item v-for="item in val" :key="item.code">
                    <n-grid :cols="3" :y-gap="0">
                      <n-gi>

                        <n-text :type="item.zdf>0?'error':'success'"><n-image :src="item.img"  width="20"/> &nbsp;{{ item.name }}</n-text>
                      </n-gi>
                      <n-gi>
                        <n-text :type="item.zdf>0?'error':'success'">{{ item.zxj }}</n-text>&nbsp;
                        <n-text :type="item.zdf>0?'error':'success'"><n-number-animation :precision="2" :from="0" :to="item.zdf" />%</n-text>

                      </n-gi>
                      <n-gi>
                        <n-text :type="item.state === 'open' ? 'success' : 'warning'">{{
                            item.state === 'open' ? '开市' : '休市'
                          }}
                        </n-text>
                      </n-gi>
                    </n-grid>
                  </n-list-item>
                </n-list>
              </n-gi>
            </n-grid>
          </n-tab-pane>
          <n-tab-pane name="上证指数" tab="上证指数">
            <k-line-chart code="sh000001" :chart-height="panelHeight" name="上证指数" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="深证成指" tab="深证成指">
            <k-line-chart code="sz399001" :chart-height="panelHeight" name="深证成指" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="创业板指" tab="创业板指">
            <k-line-chart code="sz399006" :chart-height="panelHeight" name="创业板指" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="恒生指数" tab="恒生指数">
            <k-line-chart code="hkHSI" :chart-height="panelHeight" name="恒生指数" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="纳斯达克" tab="纳斯达克">
            <k-line-chart code="us.IXIC" :chart-height="panelHeight" name="纳斯达克" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="道琼斯" tab="道琼斯">
            <k-line-chart code="us.DJI" :chart-height="panelHeight" name="道琼斯" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="标普500" tab="标普500">
            <k-line-chart code="us.INX" :chart-height="panelHeight" name="标普500" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
      <n-tab-pane name="指标行情" tab="指标行情">
        <n-tabs type="segment" animated>
          <n-tab-pane name="科创50" tab="科创50">
            <k-line-chart code="sh000688" :chart-height="panelHeight" name="科创50" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="沪深300" tab="沪深300">
            <k-line-chart code="sh000300" :chart-height="panelHeight" name="沪深300" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="上证50" tab="上证50">
            <k-line-chart code="sh000016" :chart-height="panelHeight" name="上证50" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="中证A500" tab="中证A500">
            <k-line-chart code="sh000510" :chart-height="panelHeight" name="中证A500" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="中证1000" tab="中证1000">
            <k-line-chart code="sh000852" :chart-height="panelHeight" name="中证1000" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="中证白酒" tab="中证白酒">
            <k-line-chart code="sz399997" :chart-height="panelHeight" name="中证白酒" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
          <n-tab-pane name="富时中国三倍做多" tab="富时中国三倍做多">
            <k-line-chart code="usYINN.AM" :chart-height="panelHeight" name="富时中国三倍做多" :k-days="20"
                          :dark-theme="true"></k-line-chart>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
    </n-tabs>
  </n-card>
  <n-modal transform-origin="center" v-model:show="summaryModal"  preset="card" style="width: 800px;" :title="'AI市场资讯总结'" >
    <n-spin size="small" :show="loading">
      <MdPreview   style="height: 440px;text-align: left" :modelValue="aiSummary" :theme="theme"/>
    </n-spin>
    <template #footer>
      <n-flex justify="space-between" ref="tipsRef">
        <n-text type="info" v-if="aiSummaryTime" >
          <n-tag v-if="modelName" type="warning" round :title="chatId" :bordered="false">{{modelName}}</n-tag>
          {{aiSummaryTime}}
        </n-text>
        <n-text type="error" >*AI分析结果仅供参考，请以实际行情为准。投资需谨慎，风险自担。</n-text>
      </n-flex>
    </template>
    <template #action>
      <n-flex justify="space-between" style="margin-bottom: 10px">
        <n-select style="width: 49%" v-model:value="sysPromptId" label-field="name" value-field="ID" :options="sysPromptOptions" placeholder="请选择系统提示词" />
        <n-select style="width: 49%" v-model:value="question" label-field="name" value-field="content" :options="userPromptOptions" placeholder="请选择用户提示词" />
      </n-flex>
      <n-flex justify="right">
        <n-input v-model:value="question" style="text-align: left"  clearable
                 type="textarea"
                 :show-count="true"
                 placeholder="请输入您的问题:例如 总结和分析股票市场新闻中的投资机会"
                 :autosize="{
              minRows: 2,
              maxRows: 5
            }"
        />
        <n-button size="tiny"  type="warning" @click="reAiSummary">再次总结</n-button>
        <n-button size="tiny" type="success" @click="copyToClipboard">复制到剪切板</n-button>
        <n-button size="tiny" type="primary" @click="saveAsMarkdown">保存为Markdown文件</n-button>
        <n-button size="tiny" type="error" @click="share">分享到项目社区</n-button>
      </n-flex>
    </template>
  </n-modal>

  <div style="position: fixed;bottom: 18px;right:25px;z-index: 10;" v-if="summaryBTN">
    <n-input-group >
      <n-button type="primary" @click="getAiSummary">
        <n-icon :component="PulseOutline"/> &nbsp;AI总结
      </n-button>
    </n-input-group>
  </div>

</template>
<style scoped>

</style>