<template>
  <div class="template-page">
    <header class="template-header">
      <div>
        <button class="back-link" type="button" @click="router.push('/app/resumes')">
          <ArrowLeftOutlined /> 返回简历实验室
        </button>
        <p class="step-label">创建简历 · 第 1 步</p>
        <h1>先选择一张适合你的模板</h1>
        <p class="header-copy">模板只改变排版，不锁定内容。进入编辑器后，所有栏目都可以自由填写、隐藏和实时预览。</p>
      </div>
      <div class="header-note">
        <SafetyCertificateOutlined />
        <span><strong>文本化模板</strong><small>没有图片切片，导出更清晰，也更利于 ATS 识别</small></span>
      </div>
    </header>

    <div class="filter-row" aria-label="模板说明">
      <span>首批 10 款</span>
      <i />
      <span>全部免费</span>
      <i />
      <span>支持实时切换</span>
    </div>

    <main class="template-grid" aria-label="简历模板列表">
      <button
        v-for="(template, index) in resumeTemplates"
        :key="template.id"
        class="template-card"
        :class="{ creating: creatingTemplate === template.id }"
        type="button"
        :aria-label="`使用${template.name}模板`"
        :disabled="Boolean(creatingTemplate)"
        @click="createFromTemplate(template)"
      >
        <div class="template-sequence">{{ String(index + 1).padStart(2, '0') }}</div>
        <div class="paper-frame">
          <div class="paper-preview" :class="`preview-${template.preview}`" :style="{ '--template-accent': template.accent }">
            <div class="preview-sidebar" />
            <div class="preview-content">
              <div class="preview-name">林之遥</div>
              <div class="preview-role">PRODUCT DESIGNER</div>
              <div class="preview-contact">Shanghai · hello@example.com</div>
              <div class="preview-rule" />
              <section v-for="section in 3" :key="section">
                <b>{{ section === 1 ? 'EXPERIENCE' : section === 2 ? 'PROJECTS' : 'EDUCATION' }}</b>
                <span />
                <span class="short" />
                <span />
              </section>
            </div>
          </div>
          <div class="use-overlay">
            <span v-if="creatingTemplate === template.id"><LoadingOutlined spin /> 正在创建</span>
            <span v-else>使用此模板 <ArrowRightOutlined /></span>
          </div>
        </div>
        <div class="template-meta">
          <div><strong>{{ template.name }}</strong><span>{{ template.tone }}</span></div>
          <p>{{ template.description }}</p>
          <small>{{ template.audience }}</small>
        </div>
      </button>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  ArrowLeftOutlined,
  ArrowRightOutlined,
  LoadingOutlined,
  SafetyCertificateOutlined,
} from '@ant-design/icons-vue'
import { createTemplateContent, resumeTemplates, type ResumeTemplate } from '@/data/resumeTemplates'
import { useResumeStore } from '@/stores/resume'

const router = useRouter()
const resumeStore = useResumeStore()
const creatingTemplate = ref<string | null>(null)

const createFromTemplate = async (template: ResumeTemplate) => {
  creatingTemplate.value = template.id
  const created = await resumeStore.create({
    name: `${template.name}简历`,
    scene: 'manual',
    initial_content: createTemplateContent(template.id),
  })
  creatingTemplate.value = null

  if (!created) {
    message.error('模板创建失败，请稍后重试')
    return
  }
  router.push(`/app/resumes/${created.id}`)
}
</script>

<style scoped>
.template-page {
  --ink: #17211d;
  --muted: #6e7973;
  --green: #087443;
  width: min(1480px, 100%);
  margin: 0 auto;
  color: var(--ink);
}
.template-header {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  align-items: end;
  gap: 48px;
  padding: 12px 0 32px;
  border-bottom: 1px solid #dfe5e1;
}
.back-link {
  margin: 0 0 28px;
  padding: 0;
  border: 0;
  background: transparent;
  color: #5e6a64;
  cursor: pointer;
  font-size: 13px;
}
.back-link:hover { color: var(--green); }
.step-label {
  margin: 0 0 8px;
  color: var(--green);
  font: 700 11px/1.2 ui-monospace, SFMono-Regular, Menlo, monospace;
  letter-spacing: .12em;
}
.template-header h1 {
  margin: 0;
  font-family: 'Songti SC', 'Noto Serif SC', serif;
  font-size: clamp(34px, 4vw, 54px);
  font-weight: 650;
  letter-spacing: -.045em;
  line-height: 1.08;
}
.header-copy {
  max-width: 680px;
  margin: 15px 0 0;
  color: var(--muted);
  font-size: 14px;
  line-height: 1.75;
}
.header-note {
  display: flex;
  gap: 12px;
  padding: 18px 20px;
  border-left: 3px solid var(--green);
  background: #f1f6f3;
}
.header-note > :deep(.anticon) { margin-top: 2px; color: var(--green); font-size: 19px; }
.header-note span { display: flex; flex-direction: column; gap: 4px; }
.header-note strong { font-size: 13px; }
.header-note small { color: var(--muted); font-size: 11px; line-height: 1.5; }
.filter-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px 0;
  color: #748079;
  font-size: 11px;
}
.filter-row span:first-child { color: var(--ink); font-weight: 750; }
.filter-row i { width: 3px; height: 3px; border-radius: 50%; background: #aab3ae; }
.template-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(170px, 1fr));
  gap: 30px 22px;
  padding-bottom: 42px;
}
.template-card {
  min-width: 0;
  position: relative;
  padding: 0;
  border: 0;
  background: transparent;
  color: inherit;
  text-align: left;
  cursor: pointer;
}
.template-card:disabled { cursor: wait; }
.template-sequence {
  position: absolute;
  top: 10px;
  left: -10px;
  z-index: 3;
  color: #7d8982;
  font: 700 9px/1 ui-monospace, SFMono-Regular, Menlo, monospace;
  writing-mode: vertical-rl;
}
.paper-frame {
  position: relative;
  aspect-ratio: 210 / 297;
  margin-left: 9px;
  overflow: hidden;
  border: 1px solid #d9dfdb;
  background: #fff;
  box-shadow: 0 12px 30px rgba(36, 48, 41, .08);
  transition: transform .24s ease, box-shadow .24s ease, border-color .24s ease;
}
.template-card:hover .paper-frame,
.template-card:focus-visible .paper-frame {
  transform: translateY(-5px);
  border-color: #7da690;
  box-shadow: 0 20px 38px rgba(31, 57, 43, .15);
}
.template-card:focus-visible { outline: 2px solid var(--green); outline-offset: 6px; }
.paper-preview {
  width: 100%;
  height: 100%;
  position: relative;
  padding: 14% 12%;
  overflow: hidden;
  background: #fff;
  color: #1b211e;
  font-family: 'PingFang SC', sans-serif;
}
.preview-content { position: relative; z-index: 1; }
.preview-name { font-size: clamp(12px, 1.4vw, 21px); font-weight: 820; letter-spacing: .04em; }
.preview-role { margin-top: 4px; color: var(--template-accent); font-size: 5px; font-weight: 800; letter-spacing: .15em; }
.preview-contact { margin-top: 6px; color: #747c78; font-size: 4px; }
.preview-rule { height: 2px; margin: 9px 0 10px; background: var(--template-accent); }
.paper-preview section { margin-top: 9px; }
.paper-preview section b { display: block; margin-bottom: 5px; color: var(--template-accent); font-size: 5px; letter-spacing: .08em; }
.paper-preview section span { display: block; width: 100%; height: 2px; margin-top: 4px; border-radius: 2px; background: #cbd1cd; }
.paper-preview section span.short { width: 68%; }
.preview-band { padding-top: 20%; background: linear-gradient(#e9f3ee 0 25%, #fff 25%); }
.preview-formal { font-family: Georgia, 'Songti SC', serif; }
.preview-formal .preview-name { color: #25354a; text-align: center; }
.preview-formal .preview-role, .preview-formal .preview-contact { text-align: center; }
.preview-dense { padding: 9% 10%; }
.preview-dense section { margin-top: 6px; }
.preview-split { padding-left: 39%; background: linear-gradient(90deg, #214c3d 0 31%, #fff 31%); }
.preview-split::before { content: ''; width: 13%; aspect-ratio: 1; position: absolute; top: 12%; left: 9%; border: 1px solid rgba(255,255,255,.7); border-radius: 50%; }
.preview-magazine .preview-name { max-width: 80%; color: #b5523b; font-family: Georgia, 'Songti SC', serif; font-size: clamp(16px, 2vw, 28px); line-height: .9; }
.preview-magazine .preview-rule { width: 28%; }
.preview-air { padding: 19% 14%; }
.preview-air .preview-rule { height: 1px; background: #9ba39e; }
.preview-paper { padding: 12% 10%; font-family: Georgia, 'Songti SC', serif; }
.preview-paper .preview-name, .preview-paper .preview-role, .preview-paper .preview-contact { text-align: center; }
.preview-paper .preview-rule { height: 1px; background: #4f514f; }
.preview-block { padding-top: 28%; background: linear-gradient(#126b51 0 30%, #fff 30%); }
.preview-block .preview-name, .preview-block .preview-role, .preview-block .preview-contact { color: #fff; }
.preview-fresh { background: linear-gradient(135deg, #eef8f7 0 20%, #fff 20%); }
.preview-fresh .preview-rule { border-radius: 8px; background: #277c87; }
.use-overlay {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  background: rgba(18, 35, 27, .72);
  color: #fff;
  opacity: 0;
  transition: opacity .2s ease;
  backdrop-filter: blur(2px);
}
.use-overlay span { display: flex; align-items: center; gap: 8px; font-size: 12px; font-weight: 750; }
.template-card:hover .use-overlay,
.template-card:focus-visible .use-overlay,
.template-card.creating .use-overlay { opacity: 1; }
.template-meta { padding: 13px 4px 0 10px; }
.template-meta > div { display: flex; align-items: center; gap: 8px; }
.template-meta strong { font-size: 14px; }
.template-meta span { padding: 2px 6px; border: 1px solid #d7ded9; border-radius: 2px; color: #748079; font-size: 9px; }
.template-meta p { min-height: 39px; margin: 7px 0 6px; color: var(--muted); font-size: 10.5px; line-height: 1.55; }
.template-meta small { color: var(--green); font-size: 9.5px; font-weight: 700; }
@media (max-width: 1250px) { .template-grid { grid-template-columns: repeat(4, minmax(170px, 1fr)); } }
@media (max-width: 980px) { .template-header { grid-template-columns: 1fr; gap: 22px; } .header-note { max-width: 520px; } .template-grid { grid-template-columns: repeat(3, minmax(160px, 1fr)); } }
@media (max-width: 680px) { .template-grid { grid-template-columns: repeat(2, minmax(140px, 1fr)); gap: 24px 14px; } .template-header h1 { font-size: 34px; } .header-note { display: none; } }
@media (max-width: 420px) { .template-grid { grid-template-columns: 1fr; } .paper-frame { max-width: 250px; margin-inline: auto; } }
</style>
