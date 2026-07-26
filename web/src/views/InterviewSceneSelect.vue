<template>
  <div class="scene-page">
    <header class="scene-header">
      <button class="back-button" @click="router.push('/app/interviews')">
        <ArrowLeftOutlined /> 返回训练记录
      </button>
      <div>
        <p class="step-label">面试训练场 · 选择场景</p>
        <h1>先进入一间真实的面试现场</h1>
        <p class="intro">场地、考官和流程共同决定训练内容。第一阶段开放模拟教室，其他场景将陆续加入。</p>
      </div>
      <div class="scene-index"><strong>01</strong><span>/ 12 场景已开放</span></div>
    </header>

    <main class="scene-layout">
      <section class="scene-library" aria-label="面试场景列表">
        <button
          v-for="scene in scenes"
          :key="scene.key"
          :class="['scene-card', { selected: selectedScene === scene.key, disabled: !scene.available }]"
          :style="{ backgroundPosition: scene.position }"
          :disabled="!scene.available"
          @click="selectedScene = scene.key"
        >
          <span class="scene-shade"></span>
          <span class="scene-number">{{ scene.no }}</span>
          <span class="scene-copy"><strong>{{ scene.title }}</strong><small>{{ scene.subtitle }}</small></span>
          <span :class="['scene-status', { muted: !scene.available }]">{{ scene.available ? '已开放' : '筹备中' }}</span>
        </button>
      </section>

      <section class="classroom-panel">
        <div class="classroom-preview">
          <div class="room-light"></div>
          <div class="examiner examiner-left"><span>考官</span></div>
          <div class="examiner examiner-main"><span>主考官</span></div>
          <div class="examiner examiner-right"><span>记录员</span></div>
          <div class="blackboard">
            <span class="board-label">模拟教室</span>
            <strong>{{ form.subject }} · {{ form.grade }}</strong>
            <p>结构化问答 → 限时试讲 → 答辩追问</p>
            <i></i>
          </div>
          <div class="desk"><span>考生席</span></div>
          <div class="room-caption">
            <span class="live-dot"></span>
            教资面试 · 全流程演练
          </div>
        </div>

        <div class="config-panel">
          <div class="config-heading">
            <div>
              <p>模拟教室</p>
              <h2>配置你的教资面试</h2>
            </div>
            <span>约 15 分钟</span>
          </div>

          <div class="field-grid">
            <label>
              <span>面试类型</span>
              <a-select v-model:value="form.interviewType">
                <a-select-option value="教师资格证面试">教师资格证面试</a-select-option>
                <a-select-option value="教师招聘试讲">教师招聘试讲</a-select-option>
                <a-select-option value="说课训练">说课训练</a-select-option>
              </a-select>
            </label>
            <label>
              <span>学段</span>
              <a-select v-model:value="form.grade">
                <a-select-option value="小学">小学</a-select-option>
                <a-select-option value="初中">初中</a-select-option>
                <a-select-option value="高中">高中</a-select-option>
              </a-select>
            </label>
            <label>
              <span>学科</span>
              <a-select v-model:value="form.subject">
                <a-select-option v-for="subject in subjects" :key="subject" :value="subject">{{ subject }}</a-select-option>
              </a-select>
            </label>
            <label>
              <span>考官风格</span>
              <a-select v-model:value="form.examinerStyle">
                <a-select-option value="标准规范">标准规范</a-select-option>
                <a-select-option value="温和引导">温和引导</a-select-option>
                <a-select-option value="连续追问">连续追问</a-select-option>
              </a-select>
            </label>
          </div>

          <label class="topic-field">
            <span>试讲主题（可选）</span>
            <a-input v-model:value="form.topic" placeholder="例如：小学语文《荷花》第二课时" />
          </label>

          <div class="process-strip">
            <span><b>01</b> 结构化问答</span>
            <span><b>02</b> 模拟试讲</span>
            <span><b>03</b> 考官答辩</span>
          </div>

          <button class="enter-button" :disabled="creating" @click="createTeachingInterview">
            <LoadingOutlined v-if="creating" />
            <PlayCircleOutlined v-else />
            {{ creating ? '正在布置考场' : '进入模拟教室' }}
          </button>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeftOutlined, PlayCircleOutlined, LoadingOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useInterviewStore } from '@/stores/interview'

const router = useRouter()
const interviewStore = useInterviewStore()
const selectedScene = ref('teaching')
const creating = ref(false)
const subjects = ['语文', '数学', '英语', '物理', '化学', '生物', '历史', '地理', '政治']
const form = reactive({
  interviewType: '教师资格证面试',
  grade: '小学',
  subject: '语文',
  examinerStyle: '标准规范',
  topic: '',
})

const scenes = [
  { key: 'teaching', no: '01', title: '模拟教室', subtitle: '教资面试 · 试讲 · 答辩', position: '0% 0%', available: true },
  { key: 'corporate', no: '02', title: '企业会议室', subtitle: 'HR · 业务面 · 校招', position: '33.333% 0%', available: false },
  { key: 'group', no: '03', title: '群面讨论室', subtitle: '无领导小组讨论', position: '66.666% 0%', available: false },
  { key: 'defense', no: '04', title: '项目答辩室', subtitle: '论文 · 项目 · 方案', position: '100% 0%', available: false },
  { key: 'client', no: '05', title: '客户会议室', subtitle: '销售 · 售前 · 咨询', position: '0% 50%', available: false },
  { key: 'pressure', no: '06', title: '压力面试室', subtitle: '突发追问 · 临场反应', position: '33.333% 50%', available: false },
  { key: 'public', no: '07', title: '结构化面试厅', subtitle: '公务员 · 事业单位', position: '66.666% 50%', available: false },
  { key: 'medical', no: '08', title: '医疗面试室', subtitle: '医护 · 规培 · 医患沟通', position: '100% 50%', available: false },
  { key: 'media', no: '09', title: '媒体演播室', subtitle: '主持 · 公关 · 镜头表达', position: '0% 100%', available: false },
  { key: 'remote', no: '10', title: '远程面试间', subtitle: '视频面试 · 英文面试', position: '33.333% 100%', available: false },
  { key: 'system', no: '11', title: '系统设计室', subtitle: '架构白板 · 技术评审', position: '66.666% 100%', available: false },
  { key: 'aviation', no: '12', title: '航空面试厅', subtitle: '空乘 · 服务 · 仪态', position: '100% 100%', available: false },
]

const createTeachingInterview = async () => {
  creating.value = true
  const topic = form.topic.trim() || '由考官根据学段和学科现场抽题'
  const interview = await interviewStore.create({
    scene: 'teaching',
    target_company: form.interviewType,
    target_position: `${form.grade}${form.subject}教师`,
    target_jd: `面试类型：${form.interviewType}\n学段：${form.grade}\n学科：${form.subject}\n试讲主题：${topic}\n考官风格：${form.examinerStyle}\n流程：结构化问答、模拟试讲、考官答辩。`,
    difficulty: form.examinerStyle === '连续追问' ? 'senior' : 'mid',
    total_questions: 5,
    mode: 'hybrid',
  })
  creating.value = false
  if (!interview) {
    message.error('考场创建失败，请检查大模型配置')
    return
  }
  router.push(`/app/interviews/${interview.id}`)
}
</script>

<style scoped>
.scene-page{min-height:calc(100dvh - 64px);padding:32px clamp(20px,4vw,64px) 48px;background:#edf0ea;color:#16342d}
.scene-header{max-width:1440px;margin:0 auto 28px;display:grid;grid-template-columns:180px 1fr auto;align-items:end;gap:32px}
.back-button{align-self:start;border:0;background:transparent;color:#557068;cursor:pointer;padding:8px 0;text-align:left}
.step-label{margin:0 0 8px;color:#b45c36;font-size:12px;font-weight:700;letter-spacing:.12em}
h1{margin:0;font-family:"Songti SC","STSong",serif;font-size:clamp(34px,4vw,58px);font-weight:600;letter-spacing:-.04em}
.intro{max-width:640px;margin:10px 0 0;color:#61716c;line-height:1.7}
.scene-index{display:flex;align-items:baseline;gap:8px;color:#71827c}.scene-index strong{font-size:34px;color:#173e34}.scene-index span{font-size:12px}
.scene-layout{max-width:1440px;margin:auto}
.scene-library{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:12px;margin-bottom:20px}
.scene-card{position:relative;isolation:isolate;min-height:164px;overflow:hidden;border:1px solid rgba(255,255,255,.4);background-image:url('/scenes/interview-scenes-atlas.jpg');background-size:400% 300%;display:grid;grid-template-columns:1fr auto;align-content:end;gap:4px 12px;padding:18px;text-align:left;color:#fff;cursor:pointer;box-shadow:0 8px 24px rgba(24,48,41,.09);transition:transform .22s ease,box-shadow .22s ease}
.scene-shade{position:absolute;z-index:-1;inset:0;background:linear-gradient(180deg,transparent 28%,rgba(9,25,20,.82) 100%)}
.scene-card:hover:not(.disabled),.scene-card.selected{transform:translateY(-3px);box-shadow:0 14px 30px rgba(24,48,41,.18)}.scene-card.selected{outline:3px solid #b45c36;outline-offset:-3px}.scene-card.disabled{cursor:not-allowed;filter:saturate(.72)}
.scene-number{position:absolute;left:14px;top:12px;padding:5px 7px;background:rgba(12,35,28,.76);font-size:11px}.scene-status{align-self:end;padding:4px 7px;background:#b45c36;font-size:10px}.scene-status.muted{background:rgba(26,44,38,.72);color:#d8dfdb}.scene-copy{display:flex;flex-direction:column;gap:3px}.scene-copy strong{font-size:15px}.scene-copy small{color:#dce5df}
.classroom-panel{min-width:0;background:#f8f8f3;border:1px solid #cbd2cc;display:grid;grid-template-columns:minmax(430px,1.15fr) minmax(330px,.85fr);box-shadow:0 22px 60px rgba(28,56,47,.09)}
.classroom-preview{position:relative;min-height:610px;overflow:hidden;background:linear-gradient(#dce5df 0 58%,#b38a62 58% 60%,#caa77e 60%);perspective:800px}
.room-light{position:absolute;inset:0;background:linear-gradient(115deg,rgba(255,255,255,.68),transparent 36%),repeating-linear-gradient(90deg,transparent 0 18%,rgba(255,255,255,.12) 18% 18.5%)}
.blackboard{position:absolute;left:12%;right:12%;top:18%;height:36%;padding:28px 34px;background:#173e34;color:#eef1df;border:12px solid #8d6948;box-shadow:0 12px 28px rgba(22,46,38,.25);transform:rotateX(-1deg)}
.board-label{display:block;color:#b8c9bd;font-size:11px;letter-spacing:.16em}.blackboard strong{display:block;margin-top:22px;font-family:"Songti SC",serif;font-size:28px}.blackboard p{margin-top:12px;color:#c8d6cc}.blackboard i{display:block;width:58%;height:1px;margin-top:22px;background:#90a99b}
.examiner{position:absolute;top:5%;width:88px;height:88px;border-radius:50%;background:#244f43;color:#fff;display:flex;align-items:end;justify-content:center;box-shadow:0 8px 18px rgba(14,39,32,.18)}.examiner span{transform:translateY(25px);color:#476158;font-size:12px}.examiner-left{left:23%}.examiner-main{left:44%;background:#b35d37}.examiner-right{right:23%}
.desk{position:absolute;left:28%;right:28%;bottom:10%;height:58px;background:#7b5437;color:#eadcc9;display:flex;align-items:center;justify-content:center;box-shadow:0 20px 0 #5f402c}.room-caption{position:absolute;left:22px;bottom:20px;padding:9px 12px;background:rgba(19,48,40,.9);color:#fff;font-size:12px}.live-dot{display:inline-block;width:7px;height:7px;margin-right:7px;border-radius:50%;background:#e46d43}
.config-panel{padding:34px;display:flex;flex-direction:column}.config-heading{display:flex;justify-content:space-between;gap:20px;border-bottom:1px solid #d8ddd8;padding-bottom:22px}.config-heading p{margin:0 0 4px;color:#b45c36;font-size:12px;font-weight:700}.config-heading h2{margin:0;font-family:"Songti SC",serif;font-size:28px}.config-heading>span{color:#6d7c76;font-size:12px}
.field-grid{display:grid;grid-template-columns:1fr 1fr;gap:18px 14px;margin-top:24px}.field-grid label,.topic-field{display:flex;flex-direction:column;gap:7px}.field-grid label>span,.topic-field>span{font-size:12px;color:#5c6e67}.field-grid :deep(.ant-select){width:100%}.topic-field{margin-top:18px}
.process-strip{display:grid;grid-template-columns:repeat(3,1fr);gap:8px;margin:24px 0}.process-strip span{padding:12px 8px;border-top:1px solid #b9c4be;color:#52675f;font-size:11px}.process-strip b{display:block;margin-bottom:5px;color:#b45c36}
.enter-button{margin-top:auto;border:0;background:#173e34;color:#fff;padding:15px 20px;font-weight:700;cursor:pointer;transition:.2s}.enter-button:hover{background:#245849;transform:translateY(-1px)}.enter-button:disabled{opacity:.6}
@media(max-width:1100px){.scene-header{grid-template-columns:1fr}.back-button{order:-1}.scene-index{display:none}.scene-library{grid-template-columns:repeat(3,1fr)}.classroom-panel{grid-template-columns:1fr 380px}.classroom-preview{min-height:520px}}
@media(max-width:760px){.scene-page{padding:22px 14px 36px}.scene-library{grid-template-columns:1fr 1fr}.scene-card{min-height:132px;padding:14px}.scene-status{display:none}.classroom-panel{grid-template-columns:1fr}.classroom-preview{min-height:400px}.config-panel{padding:24px}.field-grid{grid-template-columns:1fr}h1{font-size:36px}}
</style>
