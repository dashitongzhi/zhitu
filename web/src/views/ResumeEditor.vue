<template>
  <div class="resume-shell" :class="{ 'has-preview-sidebar': demoMode }">
    <aside v-if="demoMode" class="product-sidebar" aria-label="产品导航">
      <button class="sidebar-brand" type="button" @click="router.push('/')">
        <span class="sidebar-logo">职</span>
        <span><strong>职途</strong><small>智能求职工作台</small></span>
      </button>
      <nav class="product-nav">
        <button class="active" type="button" aria-current="page">
          <FileTextOutlined />
          <span><strong>简历实验室</strong><small>编辑与实时预览</small></span>
        </button>
        <button type="button" @click="previewNavigate('interview')">
          <CommentOutlined />
          <span><strong>面试训练场</strong><small>模拟面试与复盘</small></span>
        </button>
        <button type="button" @click="previewNavigate('delivery')">
          <SendOutlined />
          <span><strong>投递看板</strong><small>跟踪求职进度</small></span>
        </button>
      </nav>
      <div class="sidebar-foot">
        <span class="status-dot" />
        本地设计预览
      </div>
    </aside>

    <div class="resume-lab" :class="{ 'is-demo': demoMode }">
    <header class="lab-toolbar">
      <div class="toolbar-leading">
        <button class="icon-button" type="button" aria-label="返回简历列表" @click="backToList">
          <ArrowLeftOutlined />
        </button>
        <div class="brand-mark"><span>职</span></div>
        <div class="resume-title-wrap">
          <input v-model="editableName" class="resume-name" @blur="handleNameSave" />
          <span class="save-state"><CheckCircleOutlined /> {{ demoMode ? '本地预览模式' : '内容实时预览' }}</span>
        </div>
      </div>

      <div class="toolbar-controls">
        <label class="toolbar-select">
          <LayoutOutlined />
          <select v-model="templateStyle">
            <option v-for="template in resumeTemplates" :key="template.id" :value="template.id">
              {{ template.name }}
            </option>
          </select>
        </label>
        <label class="toolbar-select">
          <FontSizeOutlined />
          <select v-model="fontFamily">
            <option value="serif">宋体正文</option>
            <option value="sans">黑体正文</option>
          </select>
        </label>
        <label class="toolbar-select compact">
          <ColumnHeightOutlined />
          <select v-model="density">
            <option value="compact">紧凑</option>
            <option value="comfortable">舒展</option>
          </select>
        </label>
        <div class="zoom-control">
          <button type="button" @click="zoom = Math.max(60, zoom - 5)">−</button>
          <span>{{ zoom }}%</span>
          <button type="button" @click="zoom = Math.min(110, zoom + 5)">＋</button>
        </div>
        <a-button @click="versionDrawerOpen = true"><HistoryOutlined /> 版本</a-button>
        <a-button @click="showSaveVersionModal = true"><SaveOutlined /> 保存版本</a-button>
        <a-button class="export-button" type="primary" @click="exportResume">
          <DownloadOutlined /> 导出 PDF
        </a-button>
      </div>
    </header>

    <main class="lab-workspace">
      <aside class="editor-panel">
        <div class="panel-tabs">
          <button
            type="button"
            :class="{ active: editorTab === 'edit' }"
            @click="editorTab = 'edit'"
          >结构化编辑</button>
          <button
            type="button"
            :class="{ active: editorTab === 'smart' }"
            @click="editorTab = 'smart'"
          >
            <BulbOutlined /> 智能完善
          </button>
        </div>

        <!-- 结构化编辑 -->
        <div v-if="editorTab === 'edit'" class="editor-scroll">
          <div class="panel-hint">填写内容后，右侧 A4 简历会即时更新</div>

          <section class="form-section" :class="{ collapsed: collapsed.personal }">
            <div class="section-heading" :class="{ disabled: !isVisible('personal') }">
              <label class="visibility-toggle" title="控制基本信息是否显示在简历中">
                <input v-model="resumeContent.module_visibility.personal" type="checkbox" />
                <span class="custom-check"><CheckOutlined /></span>
                <UserOutlined /><strong>基本信息</strong>
              </label>
              <button class="collapse-button" type="button" aria-label="展开或收起基本信息" @click="toggleSection('personal')"><DownOutlined /></button>
            </div>
            <div class="section-body grid-two">
              <label class="field full"><span>姓名</span><input v-model="resumeContent.personal.name" /></label>
              <label class="field"><span>电话</span><input v-model="resumeContent.personal.phone" /></label>
              <label class="field"><span>邮箱</span><input v-model="resumeContent.personal.email" /></label>
              <label class="field"><span>所在城市</span><input v-model="resumeContent.personal.city" /></label>
              <label class="field"><span>GitHub / 主页</span><input v-model="resumeContent.personal.github" /></label>
            </div>
          </section>

          <section class="form-section" :class="{ collapsed: collapsed.intention }">
            <div class="section-heading" :class="{ disabled: !isVisible('intention') }">
              <label class="visibility-toggle" title="控制求职意向是否显示在简历中">
                <input v-model="resumeContent.module_visibility.intention" type="checkbox" />
                <span class="custom-check"><CheckOutlined /></span>
                <AimOutlined /><strong>求职意向</strong>
              </label>
              <button class="collapse-button" type="button" aria-label="展开或收起求职意向" @click="toggleSection('intention')"><DownOutlined /></button>
            </div>
            <div class="section-body grid-two">
              <label class="field full"><span>目标岗位</span><input v-model="resumeContent.intention.position" /></label>
              <label class="field"><span>期望城市</span><input v-model="resumeContent.intention.city" /></label>
              <label class="field"><span>期望薪资</span><input v-model="resumeContent.intention.salary" /></label>
              <label class="field"><span>到岗时间</span><input v-model="resumeContent.intention.arrival" /></label>
              <label class="field"><span>目标行业</span><input v-model="resumeContent.intention.industry" /></label>
            </div>
          </section>

          <section class="form-section" :class="{ collapsed: collapsed.education }">
            <div class="section-heading" :class="{ disabled: !isVisible('education') }">
              <label class="visibility-toggle" title="控制教育经历是否显示在简历中">
                <input v-model="resumeContent.module_visibility.education" type="checkbox" />
                <span class="custom-check"><CheckOutlined /></span>
                <ReadOutlined /><strong>教育经历</strong>
              </label>
              <button class="collapse-button" type="button" aria-label="展开或收起教育经历" @click="toggleSection('education')"><DownOutlined /></button>
            </div>
            <div class="section-body repeat-list">
              <article v-for="(item, index) in resumeContent.education" :key="`edu-${index}`" class="repeat-item">
                <button class="remove-item" type="button" @click="resumeContent.education.splice(index, 1)"><DeleteOutlined /></button>
                <div class="grid-two">
                  <label class="field full"><span>学校</span><input v-model="item.school" /></label>
                  <label class="field"><span>专业</span><input v-model="item.major" /></label>
                  <label class="field"><span>学历</span><input v-model="item.degree" /></label>
                  <label class="field"><span>开始</span><input v-model="item.start" /></label>
                  <label class="field"><span>结束</span><input v-model="item.end" /></label>
                  <label class="field full"><span>主修课程 / GPA</span><input v-model="item.courses" /></label>
                </div>
              </article>
              <button class="add-item" type="button" @click="addEducation"><PlusOutlined /> 添加教育经历</button>
            </div>
          </section>

          <section class="form-section" :class="{ collapsed: collapsed.work }">
            <div class="section-heading" :class="{ disabled: !isVisible('work') }">
              <label class="visibility-toggle" title="控制工作经历是否显示在简历中">
                <input v-model="resumeContent.module_visibility.work" type="checkbox" />
                <span class="custom-check"><CheckOutlined /></span>
                <BankOutlined /><strong>工作经历</strong>
              </label>
              <button class="collapse-button" type="button" aria-label="展开或收起工作经历" @click="toggleSection('work')"><DownOutlined /></button>
            </div>
            <div class="section-body repeat-list">
              <article v-for="(item, index) in resumeContent.work" :key="`work-${index}`" class="repeat-item">
                <button class="remove-item" type="button" @click="resumeContent.work.splice(index, 1)"><DeleteOutlined /></button>
                <div class="grid-two">
                  <label class="field full"><span>公司</span><input v-model="item.company" /></label>
                  <label class="field"><span>职位</span><input v-model="item.position" /></label>
                  <label class="field"><span>时间</span><input v-model="item.start" placeholder="2023.07 - 至今" /></label>
                  <label class="field full"><span>工作成果（每行一条）</span><textarea v-model="item.description" rows="5" /></label>
                </div>
              </article>
              <button class="add-item" type="button" @click="addWork"><PlusOutlined /> 添加工作经历</button>
            </div>
          </section>

          <section class="form-section" :class="{ collapsed: collapsed.project }">
            <div class="section-heading" :class="{ disabled: !isVisible('project') }">
              <label class="visibility-toggle" title="控制项目经历是否显示在简历中">
                <input v-model="resumeContent.module_visibility.project" type="checkbox" />
                <span class="custom-check"><CheckOutlined /></span>
                <ProjectOutlined /><strong>项目经历</strong>
              </label>
              <button class="collapse-button" type="button" aria-label="展开或收起项目经历" @click="toggleSection('project')"><DownOutlined /></button>
            </div>
            <div class="section-body repeat-list">
              <article v-for="(item, index) in resumeContent.project" :key="`project-${index}`" class="repeat-item">
                <button class="remove-item" type="button" @click="resumeContent.project.splice(index, 1)"><DeleteOutlined /></button>
                <div class="grid-two">
                  <label class="field full"><span>项目名称</span><input v-model="item.name" /></label>
                  <label class="field"><span>项目角色</span><input v-model="item.role" /></label>
                  <label class="field"><span>时间</span><input v-model="item.start" placeholder="2024.01 - 2024.06" /></label>
                  <label class="field full"><span>项目成果（每行一条）</span><textarea v-model="item.description" rows="5" /></label>
                  <label class="field full"><span>技术栈（逗号分隔）</span><input :value="item.tech_stack.join(', ')" @input="updateTechStack(item, $event)" /></label>
                </div>
              </article>
              <button class="add-item" type="button" @click="addProject"><PlusOutlined /> 添加项目经历</button>
            </div>
          </section>

          <section class="form-section" :class="{ collapsed: collapsed.skills }">
            <div class="section-heading" :class="{ disabled: !isVisible('skills') }">
              <label class="visibility-toggle" title="控制专业技能是否显示在简历中">
                <input v-model="resumeContent.module_visibility.skills" type="checkbox" />
                <span class="custom-check"><CheckOutlined /></span>
                <ToolOutlined /><strong>专业技能</strong>
              </label>
              <button class="collapse-button" type="button" aria-label="展开或收起专业技能" @click="toggleSection('skills')"><DownOutlined /></button>
            </div>
            <div class="section-body repeat-list">
              <article v-for="(item, index) in resumeContent.skills" :key="`skill-${index}`" class="skill-row">
                <input v-model="item.category" placeholder="分类" />
                <input v-model="item.name" placeholder="技能内容" />
                <button type="button" @click="resumeContent.skills.splice(index, 1)"><CloseOutlined /></button>
              </article>
              <button class="add-item" type="button" @click="addSkill"><PlusOutlined /> 添加技能分类</button>
            </div>
          </section>
        </div>

        <!-- 智能完善（与结构化编辑同面板切换） -->
        <div v-if="editorTab === 'smart'" class="editor-scroll smart-pane-wrap">
          <!-- 智能完善子 tabs -->
          <div class="smart-tabs">
            <button
              v-for="t in smartTabs"
              :key="t.key"
              class="smart-tab-btn"
              :class="{ active: activeSmartTab === t.key }"
              type="button"
              @click="activeSmartTab = t.key"
            >
              <component :is="t.icon" />
              <span>{{ t.label }}</span>
            </button>
          </div>

          <!-- AI 大模型分析 -->
          <div v-if="activeSmartTab === 'analysis'" class="smart-pane">
            <div class="pane-head">
              <span class="pane-title">简历内容分析</span>
              <button class="mini-btn" type="button" @click="runAnalysis" :disabled="analysisLoading">
                <LoadingOutlined v-if="analysisLoading" />
                <ReloadOutlined v-else />
                {{ analysisLoading ? '分析中…' : '重新分析' }}
              </button>
            </div>

            <div class="score-overview">
              <div class="score-ring">
                <svg viewBox="0 0 80 80">
                  <circle cx="40" cy="40" r="34" fill="none" stroke="#e8ece8" stroke-width="6" />
                  <circle
                    cx="40" cy="40" r="34" fill="none"
                    stroke="url(#scoreGrad)" stroke-width="6"
                    stroke-linecap="round"
                    :stroke-dasharray="`${(analysisScore / 100) * 213.6} 213.6`"
                    transform="rotate(-90 40 40)"
                  />
                  <defs>
                    <linearGradient id="scoreGrad" x1="0" y1="0" x2="1" y2="1">
                      <stop offset="0%" stop-color="#087443" />
                      <stop offset="100%" stop-color="#34c759" />
                    </linearGradient>
                  </defs>
                </svg>
                <div class="score-num">{{ analysisScore }}<small>/100</small></div>
              </div>
              <div class="score-meta">
                <strong>{{ analysisScoreLabel }}</strong>
                <p>{{ analysisScoreDesc }}</p>
              </div>
            </div>

            <section class="dim-card">
              <h4>多维度评分</h4>
              <div v-for="d in analysisDimensions" :key="d.label" class="dim-row">
                <span class="dim-label">{{ d.label }}</span>
                <div class="dim-track"><i :style="{ width: `${d.value}%`, background: d.color }" /></div>
                <b :style="{ color: d.color }">{{ d.value }}</b>
              </div>
            </section>

            <section class="dim-card">
              <h4>亮点</h4>
              <div v-if="analysisHighlights.length" class="hl-list">
                <article v-for="(h, i) in analysisHighlights" :key="i" class="hl-item hl-good">
                  <CheckCircleOutlined />
                  <div><strong>{{ h.title }}</strong><p>{{ h.detail }}</p></div>
                </article>
              </div>
              <a-empty v-else :image="null" description="暂无亮点" />
            </section>

            <section class="dim-card">
              <h4>待补强项</h4>
              <div v-if="analysisWeakness.length" class="hl-list">
                <article v-for="(w, i) in analysisWeakness" :key="i" class="hl-item hl-bad">
                  <ExclamationCircleOutlined />
                  <div>
                    <strong>{{ w.title }}</strong><p>{{ w.detail }}</p>
                    <button v-if="w.action" class="link-btn" type="button" @click="w.action">{{ w.actionLabel || '立即完善' }}</button>
                  </div>
                </article>
              </div>
              <a-empty v-else :image="null" description="简历已经比较完善" />
            </section>
          </div>

          <!-- JD 智能匹配 -->
          <div v-if="activeSmartTab === 'jd'" class="smart-pane">
            <div class="pane-head">
              <span class="pane-title">岗位描述匹配</span>
              <button class="mini-btn primary" type="button" @click="runJdMatch" :disabled="jdLoading">
                <LoadingOutlined v-if="jdLoading" />
                <ThunderboltOutlined v-else />
                {{ jdLoading ? '匹配中…' : '开始匹配' }}
              </button>
            </div>

            <label class="jd-textarea-wrap">
              <textarea
                v-model="targetJd"
                rows="6"
                placeholder="粘贴目标岗位 JD（职位描述、要求、加分项），AI 会自动提取关键词并计算匹配度"
              />
              <span class="jd-counter">{{ targetJd.length }} 字</span>
            </label>

            <div v-if="jdMatchResult" class="jd-result">
              <div class="jd-match-ring">
                <svg viewBox="0 0 100 100">
                  <circle cx="50" cy="50" r="42" fill="none" stroke="#e8ece8" stroke-width="8" />
                  <circle
                    cx="50" cy="50" r="42" fill="none"
                    :stroke="jdMatchColor" stroke-width="8"
                    stroke-linecap="round"
                    :stroke-dasharray="`${(jdMatchResult.matchRate / 100) * 263.9} 263.9`"
                    transform="rotate(-90 50 50)"
                  />
                </svg>
                <div class="jd-match-num">
                  <strong>{{ jdMatchResult.matchRate }}%</strong>
                  <small>{{ jdMatchLabel }}</small>
                </div>
              </div>

              <div class="jd-stats">
                <div class="jd-stat-item">
                  <span class="jd-stat-num green">{{ jdMatchResult.matched.length }}</span>
                  <span class="jd-stat-label">已匹配关键词</span>
                </div>
                <div class="jd-stat-item">
                  <span class="jd-stat-num red">{{ jdMatchResult.missing.length }}</span>
                  <span class="jd-stat-label">缺失关键词</span>
                </div>
                <div class="jd-stat-item">
                  <span class="jd-stat-num blue">{{ jdMatchResult.suggest.length }}</span>
                  <span class="jd-stat-label">建议补充</span>
                </div>
              </div>

              <div class="kw-section">
                <h5>已匹配 <small>简历中已包含</small></h5>
                <div class="kw-tags">
                  <span v-for="k in jdMatchResult.matched" :key="k" class="kw-tag kw-green">{{ k }}</span>
                  <span v-if="!jdMatchResult.matched.length" class="kw-empty">暂无</span>
                </div>
              </div>

              <div class="kw-section">
                <h5>缺失关键词 <small>JD 中出现但简历未提及</small></h5>
                <div class="kw-tags">
                  <span v-for="k in jdMatchResult.missing" :key="k" class="kw-tag kw-red">{{ k }}</span>
                  <span v-if="!jdMatchResult.missing.length" class="kw-empty">无缺失</span>
                </div>
              </div>

              <div class="kw-section">
                <h5>AI 建议 <small>基于差距分析</small></h5>
                <article v-for="(s, i) in jdMatchResult.suggest" :key="i" class="suggest-item">
                  <b>{{ i + 1 }}</b>
                  <div><strong>{{ s.title }}</strong><p>{{ s.detail }}</p></div>
                </article>
                <a-empty v-if="!jdMatchResult.suggest.length" :image="null" description="匹配良好，无补充建议" />
              </div>
            </div>

            <div v-else class="jd-empty">
              <FileSearchOutlined />
              <p>粘贴 JD 后点击「开始匹配」</p>
            </div>
          </div>

          <!-- AI 一键优化 -->
          <div v-if="activeSmartTab === 'optimize'" class="smart-pane">
            <div class="pane-head">
              <span class="pane-title">AI 内容优化</span>
            </div>

            <div class="opt-section">
              <h4>批量优化</h4>
              <p class="opt-desc">点击下方按钮，AI 会根据内容自动强化措辞、补充量化指标、提升专业度</p>
              <div class="opt-actions">
                <button class="opt-btn" type="button" @click="optimizeAll" :disabled="optimizeLoading">
                  <LoadingOutlined v-if="optimizeLoading" />
                  <ThunderboltOutlined v-else />
                  {{ optimizeLoading ? 'AI 优化中…' : '一键优化全部内容' }}
                </button>
              </div>
            </div>

            <div class="opt-section">
              <h4>逐项优化</h4>
              <article v-for="item in optimizeItems" :key="item.key" class="opt-item">
                <div class="opt-item-head">
                  <strong>{{ item.title }}</strong>
                  <span class="opt-item-status" :class="item.statusClass">{{ item.statusText }}</span>
                </div>
                <p class="opt-item-desc">{{ item.desc }}</p>
                <button class="link-btn" type="button" @click="item.run" :disabled="item.loading">
                  <LoadingOutlined v-if="item.loading" />
                  {{ item.loading ? '处理中…' : item.buttonText }}
                </button>
              </article>
            </div>
          </div>

          <p class="privacy-note"><SafetyCertificateOutlined /> 简历内容仅用于当前编辑与诊断，不会上传至第三方</p>
        </div>
      </aside>

      <section class="preview-stage">
        <div class="stage-meta">
          <span><FileTextOutlined /> A4 · 实时预览</span>
          <span>{{ contentStats.chars }} 字符 · {{ contentStats.sections }} 个内容模块</span>
        </div>
        <div class="paper-viewport">
          <article
            id="resume-paper"
            class="resume-paper"
            :class="[`template-${templateStyle}`, `font-${fontFamily}`, `density-${density}`]"
            :style="{ zoom: zoom / 100 }"
          >
            <header v-if="isVisible('personal')" class="paper-header">
              <h1>{{ resumeContent.personal.name || '你的姓名' }}</h1>
              <p class="paper-role">{{ resumeContent.intention.position || '目标职位' }}</p>
              <div class="contact-line">
                <span v-if="resumeContent.personal.phone">{{ resumeContent.personal.phone }}</span>
                <span v-if="resumeContent.personal.email">{{ resumeContent.personal.email }}</span>
                <span v-if="resumeContent.personal.city">{{ resumeContent.personal.city }}</span>
                <span v-if="resumeContent.personal.github">{{ resumeContent.personal.github }}</span>
              </div>
            </header>

            <section v-if="isVisible('intention')" class="paper-section intention-section">
              <h2>求职意向</h2>
              <p>{{ intentionSummary }}</p>
            </section>

            <section v-if="isVisible('education') && resumeContent.education.length" class="paper-section">
              <h2>教育经历</h2>
              <article v-for="(item, index) in resumeContent.education" :key="`preview-edu-${index}`" class="resume-entry">
                <div class="entry-heading">
                  <strong>{{ item.school || '学校名称' }}</strong>
                  <span>{{ item.start }}<template v-if="item.end"> — {{ item.end }}</template></span>
                </div>
                <div class="entry-subheading"><span>{{ item.major }}</span><span>{{ item.degree }}</span></div>
                <p v-if="item.courses" class="entry-note">{{ item.courses }}</p>
              </article>
            </section>

            <section v-if="isVisible('work') && resumeContent.work.length" class="paper-section">
              <h2>工作经历</h2>
              <article v-for="(item, index) in resumeContent.work" :key="`preview-work-${index}`" class="resume-entry">
                <div class="entry-heading"><strong>{{ item.company || '公司名称' }}</strong><span>{{ item.start }}<template v-if="item.end"> — {{ item.end }}</template></span></div>
                <div class="entry-subheading"><span>{{ item.position }}</span></div>
                <ul><li v-for="(line, lineIndex) in toLines(item.description)" :key="lineIndex">{{ line }}</li></ul>
              </article>
            </section>

            <section v-if="isVisible('project') && resumeContent.project.length" class="paper-section">
              <h2>项目经历</h2>
              <article v-for="(item, index) in resumeContent.project" :key="`preview-project-${index}`" class="resume-entry">
                <div class="entry-heading"><strong>{{ item.name || '项目名称' }}</strong><span>{{ item.start }}<template v-if="item.end"> — {{ item.end }}</template></span></div>
                <div class="entry-subheading"><span>{{ item.role }}</span><span class="tech-stack">{{ item.tech_stack.join(' · ') }}</span></div>
                <ul><li v-for="(line, lineIndex) in toLines(item.description)" :key="lineIndex">{{ line }}</li></ul>
              </article>
            </section>

            <section v-if="isVisible('skills') && resumeContent.skills.length" class="paper-section skill-section">
              <h2>专业技能</h2>
              <p v-for="(item, index) in resumeContent.skills" :key="`preview-skill-${index}`"><strong>{{ item.category }}：</strong>{{ item.name }}</p>
            </section>
          </article>
        </div>
      </section>
    </main>

    <a-drawer v-model:open="versionDrawerOpen" title="版本历史" width="380">
      <a-empty v-if="!resumeStore.versions.length" description="暂无已保存版本" />
      <div v-else class="version-list">
        <button v-for="version in resumeStore.versions" :key="version.id" type="button" :class="{ active: resumeStore.currentVersion?.id === version.id }" @click="selectVersion(version)">
          <span><strong>{{ version.version_label }}</strong><small>{{ formatVersionDate(version.created_at) }}</small></span>
          <em>{{ version.change_note || '无版本备注' }}</em>
        </button>
      </div>
    </a-drawer>

    <a-modal v-model:open="showSaveVersionModal" title="保存当前版本" :confirm-loading="savingVersion" @ok="handleSaveVersion">
      <a-form layout="vertical">
        <a-form-item label="版本备注" required><a-input v-model:value="newVersionNote" placeholder="例如：针对后端工程师 JD 完成优化" /></a-form-item>
        <a-alert message="当前结构化内容会保存为一个新的简历版本" type="info" show-icon />
      </a-form>
    </a-modal>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  AimOutlined, ArrowLeftOutlined, BankOutlined, CheckCircleOutlined, CheckOutlined, CloseOutlined,
  ColumnHeightOutlined, DeleteOutlined, DownloadOutlined, DownOutlined, FileTextOutlined,
  FontSizeOutlined, HistoryOutlined, LayoutOutlined, PlusOutlined, ProjectOutlined, CommentOutlined,
  ReadOutlined, RobotOutlined, SafetyCertificateOutlined, SaveOutlined, ToolOutlined,
  UserOutlined, SendOutlined,
  ReloadOutlined, LoadingOutlined, ThunderboltOutlined, ExclamationCircleOutlined,
  FileSearchOutlined, BulbOutlined, RocketOutlined,
} from '@ant-design/icons-vue'
import { useResumeStore } from '@/stores/resume'
import { useUserProfileStore } from '@/stores/userProfile'
import { getResumeTemplate, resumeTemplates, type ResumeTemplateId } from '@/data/resumeTemplates'
import type { ResumeVersion } from '@/types/models'

interface ResumeEducation { school: string; major: string; degree: string; start: string; end: string; courses: string; gpa: string }
interface ResumeWork { company: string; position: string; start: string; end: string; description: string; leave_reason: string }
interface ResumeProject { name: string; role: string; start: string; end: string; description: string; tech_stack: string[]; url: string }
interface ResumeSkill { category: string; name: string; proficiency: string }
interface ResumeContent {
  template_style: ResumeTemplateId
  personal: { name: string; gender: string; age: string; phone: string; email: string; github: string; avatar: string; city: string }
  intention: { position: string; city: string; salary: string; arrival: string; industry: string }
  education: ResumeEducation[]; work: ResumeWork[]; project: ResumeProject[]; skills: ResumeSkill[]
  honor: unknown[]; custom: unknown[]; module_order: string[]; module_visibility: Record<string, boolean>
}

const route = useRoute()
const router = useRouter()
const resumeStore = useResumeStore()
const userProfileStore = useUserProfileStore()
const demoMode = computed(() => route.name === 'ResumeLabPreview' || route.query.demo === '1')
const resumeId = computed(() => Number(route.params.id))
const editableName = ref('后端工程师－张明')
const templateStyle = ref<ResumeTemplateId>('classic')
const fontFamily = ref('sans')
const density = ref('comfortable')
const zoom = ref(82)
const versionDrawerOpen = ref(false)
const showSaveVersionModal = ref(false)
const savingVersion = ref(false)
const newVersionNote = ref('')
const targetJd = ref('')
const collapsed = reactive<Record<string, boolean>>({ personal: false, intention: false, education: false, work: false, project: false, skills: false })
const resumeContent = reactive<ResumeContent>(createSampleContent())

const serializedContent = computed(() => JSON.stringify(resumeContent, null, 2))
const contentStats = computed(() => ({
  chars: serializedContent.value.length,
  sections:
    (resumeContent.module_visibility.personal !== false ? 1 : 0)
    + (resumeContent.module_visibility.intention !== false ? 1 : 0)
    + (resumeContent.module_visibility.education !== false && resumeContent.education.length ? 1 : 0)
    + (resumeContent.module_visibility.work !== false && resumeContent.work.length ? 1 : 0)
    + (resumeContent.module_visibility.project !== false && resumeContent.project.length ? 1 : 0)
    + (resumeContent.module_visibility.skills !== false && resumeContent.skills.length ? 1 : 0),
}))
const intentionSummary = computed(() => [resumeContent.intention.position, resumeContent.intention.city, resumeContent.intention.salary, resumeContent.intention.arrival].filter(Boolean).join(' ｜ '))

// ==================== 智能完善：AI 大模型分析 ====================
type EditorTab = 'edit' | 'smart'
const editorTab = ref<EditorTab>('edit')

type SmartTabKey = 'analysis' | 'jd' | 'optimize'
const activeSmartTab = ref<SmartTabKey>('analysis')

const smartTabs = [
  { key: 'analysis' as const, label: 'AI 分析', icon: BulbOutlined },
  { key: 'jd' as const, label: 'JD 匹配', icon: FileSearchOutlined },
  { key: 'optimize' as const, label: 'AI 优化', icon: RocketOutlined },
]

// —— 分析相关 ——
const analysisLoading = ref(false)

// 多维度评分（基于内容实时计算 + 模拟 AI 权重）
const analysisDimensions = computed(() => {
  const personal = resumeContent.personal
  const edu = resumeContent.education
  const work = resumeContent.work
  const proj = resumeContent.project
  const skills = resumeContent.skills

  const contactComplete = [personal.name, personal.phone, personal.email, personal.city].filter(Boolean).length
  const intentionComplete = [resumeContent.intention.position, resumeContent.intention.city, resumeContent.intention.salary].filter(Boolean).length

  // 量化指标占比
  const allDesc = [...work, ...proj].map((i) => i.description).join(' ')
  const hasQuant = /[0-9]+(%|万|倍|ms|次|秒|分|倍|x|X)/.test(allDesc)
  const quantCount = (allDesc.match(/[0-9]+(%|万|倍|ms|次|秒|分|x|X)/g) || []).length

  return [
    { label: '基本信息完整度', value: Math.min(98, 40 + contactComplete * 15), color: '#087443' },
    { label: '求职意向清晰度', value: intentionComplete >= 3 ? 92 : 50 + intentionComplete * 14, color: '#007aff' },
    { label: '教育背景', value: Math.min(95, 50 + edu.length * 22), color: '#5e5ce6' },
    { label: '工作经历', value: Math.min(96, 45 + work.length * 18 + (hasQuant ? 8 : 0)), color: '#af52de' },
    { label: '项目经验', value: Math.min(94, 42 + proj.length * 22 + (quantCount >= 3 ? 12 : 0)), color: '#ff9500' },
    { label: '技能描述', value: Math.min(92, 48 + skills.length * 12), color: '#34c759' },
  ]
})

const analysisScore = computed(() => {
  const ds = analysisDimensions.value
  return Math.round(ds.reduce((a, b) => a + b.value, 0) / ds.length)
})

const analysisScoreLabel = computed(() => {
  const s = analysisScore.value
  if (s >= 85) return '优秀'
  if (s >= 70) return '良好'
  if (s >= 55) return '中等'
  return '待完善'
})

const analysisScoreDesc = computed(() => {
  const s = analysisScore.value
  if (s >= 85) return '简历内容完整且量化指标充足，仅需微调即可投递。'
  if (s >= 70) return '整体表现不错，建议补充亮点项目与量化成果。'
  if (s >= 55) return '基础信息已具备，但部分模块需进一步完善。'
  return '简历内容较少，建议从基本信息、工作经历和项目经验开始完善。'
})

// 亮点分析
const analysisHighlights = computed(() => {
  const list: { title: string; detail: string }[] = []
  const work = resumeContent.work
  const proj = resumeContent.project
  const allDesc = [...work, ...proj].map((i) => i.description).join(' ')
  const quantCount = (allDesc.match(/[0-9]+(%|万|倍|ms|次|秒|分|x|X)/g) || []).length

  if (quantCount >= 3) {
    list.push({ title: '量化成果突出', detail: `已使用 ${quantCount} 处量化指标，能有效展示业务影响力。` })
  }
  if (proj.length >= 2) {
    list.push({ title: '项目证据充分', detail: `共 ${proj.length} 个项目，能多角度展示技术能力。` })
  }
  if (work.length >= 2) {
    list.push({ title: '工作经历稳定', detail: `${work.length} 段工作经历，体现持续成长。` })
  }
  if (resumeContent.skills.length >= 4) {
    list.push({ title: '技能栈丰富', detail: `${resumeContent.skills.length} 项技能分类，覆盖面广。` })
  }
  return list
})

// 待补强项
const analysisWeakness = computed(() => {
  const list: {
    title: string
    detail: string
    actionLabel?: string
    action?: () => void
  }[] = []
  const work = resumeContent.work
  const proj = resumeContent.project
  const skills = resumeContent.skills
  const allDesc = [...work, ...proj].map((i) => i.description).join(' ')
  const hasQuant = /[0-9]+(%|万|倍|ms|次|秒|分)/.test(allDesc)

  if (!resumeContent.intention.position) {
    list.push({
      title: '求职意向缺失',
      detail: '未填写目标岗位，AI 难以做精准匹配。',
      actionLabel: '填写意向',
      action: () => { activeSmartTab.value = 'analysis'; message.info('请在左侧「求职意向」模块填写') },
    })
  }
  if (!hasQuant) {
    list.push({
      title: '缺少量化成果',
      detail: '工作与项目描述未使用数据指标，建议补充性能、规模、效率等量化结果。',
      actionLabel: 'AI 优化',
      action: () => { activeSmartTab.value = 'optimize' },
    })
  }
  if (proj.length < 2) {
    list.push({ title: '项目数量偏少', detail: `当前仅 ${proj.length} 个项目，建议补充 1-2 个能体现核心能力的项目。` })
  }
  if (skills.length < 3) {
    list.push({ title: '技能描述单薄', detail: `仅 ${skills.length} 项技能，建议按「编程语言 / 框架 / 工具」分类补全。` })
  }
  if (!resumeContent.personal.email || !resumeContent.personal.phone) {
    list.push({ title: '联系方式不全', detail: '邮箱或电话缺失，HR 无法联系到你。' })
  }
  return list
})

// 模拟 AI 分析（带 loading 动画）
const runAnalysis = async () => {
  analysisLoading.value = true
  await new Promise((r) => setTimeout(r, 800))
  analysisLoading.value = false
  message.success('AI 分析已更新')
}

// ==================== JD 智能匹配 ====================
const jdLoading = ref(false)
interface JdMatchResult {
  matchRate: number
  matched: string[]
  missing: string[]
  suggest: { title: string; detail: string }[]
}
const jdMatchResult = ref<JdMatchResult | null>(null)

// 关键词库（按技术栈分类，可扩展）
const techKeywords = [
  'Java', 'Go', 'Python', 'Node.js', 'C++', 'Rust', 'TypeScript', 'JavaScript',
  'Spring Boot', 'Spring Cloud', 'Gin', 'React', 'Vue', 'Angular', 'Next.js',
  'MySQL', 'PostgreSQL', 'Redis', 'MongoDB', 'ClickHouse', 'Elasticsearch', 'Kafka', 'RabbitMQ', 'RocketMQ',
  'Docker', 'Kubernetes', 'K8s', 'Jenkins', 'CI/CD', 'Prometheus', 'Grafana',
  '微服务', '分布式', '高并发', '高可用', '性能优化', '架构设计',
  '消息队列', '缓存', '数据库优化', '系统设计', 'DDD', '领域驱动设计',
  '机器学习', '深度学习', 'LLM', 'NLP', '推荐系统', '数据分析',
  '团队管理', '项目管理', '敏捷开发', 'Scrum', 'OKR',
  'Linux', 'Git', 'Shell', 'Python', '算法', '数据结构',
  'TOEIC', '英语', 'PMP',
]

const runJdMatch = async () => {
  if (!targetJd.value.trim()) {
    message.warning('请先粘贴 JD 内容')
    return
  }
  jdLoading.value = true
  await new Promise((r) => setTimeout(r, 900))
  // 提取 JD 中的关键词
  const jdText = targetJd.value.toLowerCase()
  const jdKeywords = techKeywords.filter((k) => jdText.includes(k.toLowerCase()))
  // 简历全文
  const resumeText = (
    resumeContent.personal.name + ' ' +
    resumeContent.intention.position + ' ' +
    resumeContent.skills.map((s) => s.name).join(' ') + ' ' +
    [...resumeContent.work, ...resumeContent.project].map((i) => i.description).join(' ')
  ).toLowerCase()
  const matched = jdKeywords.filter((k) => resumeText.includes(k.toLowerCase()))
  const missing = jdKeywords.filter((k) => !resumeText.includes(k.toLowerCase()))

  const matchRate = jdKeywords.length
    ? Math.round((matched.length / jdKeywords.length) * 100)
    : 50

  const suggest: { title: string; detail: string }[] = []
  if (missing.length > 0) {
    suggest.push({
      title: '补充技术关键词',
      detail: `简历中未提及：${missing.slice(0, 5).join('、')}。如有相关经验，建议在技能或项目描述中体现。`,
    })
  }
  if (matched.length > 0) {
    suggest.push({
      title: '突出已匹配项',
      detail: `已匹配 ${matched.length} 个关键词，建议在项目经历中给出具体应用场景与量化结果。`,
    })
  }
  if (matchRate < 60) {
    suggest.push({ title: '匹配度偏低', detail: '考虑调整简历重点，或寻找更贴合的岗位。' })
  } else if (matchRate >= 80) {
    suggest.push({ title: '匹配良好', detail: '可重点准备面试，强化项目深度与亮点。' })
  }

  jdMatchResult.value = { matchRate, matched, missing, suggest }
  jdLoading.value = false
  message.success('JD 匹配完成')
}

const jdMatchColor = computed(() => {
  const r = jdMatchResult.value?.matchRate ?? 0
  if (r >= 80) return '#34c759'
  if (r >= 60) return '#ff9500'
  return '#ff3b30'
})
const jdMatchLabel = computed(() => {
  const r = jdMatchResult.value?.matchRate ?? 0
  if (r >= 80) return '匹配良好'
  if (r >= 60) return '部分匹配'
  return '匹配度低'
})

// ==================== AI 一键优化 ====================
const optimizeLoading = ref(false)
const optStatus = reactive<Record<string, 'idle' | 'loading' | 'done'>>({
  work: 'idle', project: 'idle', skills: 'idle', summary: 'idle',
})

// 模拟 AI 优化：为经历描述补充量化指标
const enhanceDescription = (text: string): string => {
  if (!text) return text
  if (/[0-9]+(%|万|倍|ms|次|秒|分)/.test(text)) return text
  const lines = text.split('\n').filter(Boolean)
  const quantifiers = ['提升 30%', '日均 100w+ 请求', '响应时间降低 40%', '可用性达 99.9%']
  return lines.map((line, i) => `${line}，${quantifiers[i % quantifiers.length]}`).join('\n')
}

const optimizeItem = async (
  key: keyof typeof optStatus,
  apply: () => void,
) => {
  optStatus[key] = 'loading'
  await new Promise((r) => setTimeout(r, 700))
  apply()
  optStatus[key] = 'done'
  message.success('已应用 AI 优化')
}

const optimizeItems = computed(() => [
  {
    key: 'work',
    title: '工作经历强化',
    desc: 'AI 自动补充量化指标、业务影响力描述',
    buttonText: '优化工作经历',
    loading: optStatus.work === 'loading',
    statusText: optStatus.work === 'done' ? '已优化' : (optStatus.work === 'loading' ? '处理中' : '待优化'),
    statusClass: optStatus.work === 'done' ? 'st-done' : 'st-idle',
    run: () => optimizeItem('work', () => {
      resumeContent.work.forEach((w) => { w.description = enhanceDescription(w.description) })
    }),
  },
  {
    key: 'project',
    title: '项目经历强化',
    desc: '补充技术栈细节、性能指标与业务结果',
    buttonText: '优化项目经历',
    loading: optStatus.project === 'loading',
    statusText: optStatus.project === 'done' ? '已优化' : (optStatus.project === 'loading' ? '处理中' : '待优化'),
    statusClass: optStatus.project === 'done' ? 'st-done' : 'st-idle',
    run: () => optimizeItem('project', () => {
      resumeContent.project.forEach((p) => { p.description = enhanceDescription(p.description) })
    }),
  },
  {
    key: 'skills',
    title: '技能描述重组',
    desc: '按熟练度排序、补充分类与年限',
    buttonText: '优化技能描述',
    loading: optStatus.skills === 'loading',
    statusText: optStatus.skills === 'done' ? '已优化' : (optStatus.skills === 'loading' ? '处理中' : '待优化'),
    statusClass: optStatus.skills === 'done' ? 'st-done' : 'st-idle',
    run: () => optimizeItem('skills', () => {
      // 简单示例：在技能后追加「(熟练)」标记
      resumeContent.skills.forEach((s) => {
        if (!/[（(]\w+[)）]/.test(s.name)) s.name = `${s.name}（${s.proficiency || '熟练掌握'}）`
      })
    }),
  },
  {
    key: 'summary',
    title: '生成个人简介',
    desc: '基于经历自动生成一句话自我介绍',
    buttonText: 'AI 生成简介',
    loading: optStatus.summary === 'loading',
    statusText: optStatus.summary === 'done' ? '已生成' : (optStatus.summary === 'loading' ? '处理中' : '待生成'),
    statusClass: optStatus.summary === 'done' ? 'st-done' : 'st-idle',
    run: () => optimizeItem('summary', () => {
      userProfileStore.updateBasic({
        self_introduction: `${resumeContent.intention.position || '后端工程师'}，${resumeContent.work.length} 年工作经验，专注于${resumeContent.skills.slice(0, 2).map((s) => s.category).join('、')}方向，擅长高并发系统设计与性能优化。`,
      })
    }),
  },
])

const optimizeAll = async () => {
  optimizeLoading.value = true
  for (const item of optimizeItems.value) {
    await item.run()
  }
  optimizeLoading.value = false
  message.success('AI 全部优化完成')
}

watch(() => resumeStore.currentVersion, (version) => {
  if (!version || demoMode.value) return
  try {
    const normalized = normalizeContent(JSON.parse(version.content))
    Object.assign(resumeContent, normalized)
    templateStyle.value = normalized.template_style
  } catch {
    message.warning('当前版本内容不是可识别的结构化简历 JSON')
  }
})

watch(templateStyle, (value) => {
  resumeContent.template_style = value
})

function createSampleContent(): ResumeContent {
  return normalizeContent({
    personal: { name: '张明', phone: '138-0000-0000', email: 'zhangming@email.com', city: '上海市浦东新区', github: 'github.com/zhangming' },
    intention: { position: '后端工程师', city: '上海', salary: '25K–35K', arrival: '1个月内', industry: '互联网 / AI' },
    education: [{ school: '上海交通大学', major: '计算机科学与技术', degree: '本科', start: '2016.09', end: '2020.06', courses: 'GPA 3.7/4.0 · 数据结构、操作系统、数据库系统', gpa: '3.7/4.0' }],
    work: [
      { company: '上海云舟科技有限公司', position: '后端工程师', start: '2021.07', end: '至今', description: '负责公司核心业务系统设计与开发，支撑日均 1000 万+ 请求量\n拆分单体服务为 12 个微服务，系统可用性提升至 99.95%\n优化数据库慢查询与缓存策略，接口平均响应时间从 120ms 降至 45ms\n推动 CI/CD 落地，发布效率提升 70%', leave_reason: '' },
      { company: '上海智汇互联有限公司', position: '后端开发工程师', start: '2020.07', end: '2021.06', description: '参与电商平台后端开发，负责订单、支付与库存核心模块\n基于 Spring Cloud Alibaba 完成服务治理与配置管理', leave_reason: '' },
    ],
    project: [
      { name: '云舟微服务平台重构项目', role: '后端负责人', start: '2022.03', end: '2022.11', description: '设计并实现基于 Spring Cloud 的微服务框架，集成 Gateway、Sentinel、SkyWalking\n项目上线后系统稳定性提升 80%，运维成本降低 30%', tech_stack: ['Spring Cloud', 'MySQL', 'Redis'], url: '' },
      { name: '实时数据分析平台', role: '核心开发', start: '2021.09', end: '2022.02', description: '使用 Kafka + Flink 实现实时流处理，数据延迟从分钟级降至毫秒级\n基于 ClickHouse 建设查询链路，查询性能提升 5 倍', tech_stack: ['Kafka', 'Flink', 'ClickHouse'], url: '' },
    ],
    skills: [
      { category: '编程语言', name: 'Java（熟练）、Go（熟练）、Python（熟悉）', proficiency: '' },
      { category: '框架', name: 'Spring Boot、Spring Cloud Alibaba、Gin、MyBatis', proficiency: '' },
      { category: '数据库', name: 'MySQL、Redis、MongoDB、ClickHouse', proficiency: '' },
      { category: '工程工具', name: 'Docker、Kubernetes、Prometheus、Grafana', proficiency: '' },
    ],
  })
}

function normalizeContent(value: any): ResumeContent {
  return {
    template_style: getResumeTemplate(value.template_style).id,
    personal: { name: '', gender: '', age: '', phone: '', email: '', github: '', avatar: '', city: '', ...(value.personal || {}) },
    intention: { position: '', city: '', salary: '', arrival: '', industry: '', ...(value.intention || {}) },
    education: Array.isArray(value.education) ? value.education : [],
    work: Array.isArray(value.work) ? value.work : [],
    project: Array.isArray(value.project) ? value.project.map((item) => ({ ...item, tech_stack: Array.isArray(item.tech_stack) ? item.tech_stack : [] })) : [],
    skills: Array.isArray(value.skills) ? value.skills : [],
    honor: Array.isArray(value.honor) ? value.honor : [], custom: Array.isArray(value.custom) ? value.custom : [],
    module_order: value.module_order || ['personal', 'intention', 'education', 'work', 'project', 'skills', 'honor'],
    module_visibility: {
      personal: true, intention: true, education: true, work: true, project: true, skills: true, honor: true,
      ...(value.module_visibility || {}),
    },
  }
}

const toggleSection = (key: string) => { collapsed[key] = !collapsed[key] }
const isVisible = (key: string) => resumeContent.module_visibility[key] !== false
const toLines = (text: string) => text.split(/\n+/).map((line) => line.replace(/^[•·\-]\s*/, '').trim()).filter(Boolean)
const addEducation = () => resumeContent.education.push({ school: '', major: '', degree: '', start: '', end: '', courses: '', gpa: '' })
const addWork = () => resumeContent.work.push({ company: '', position: '', start: '', end: '', description: '', leave_reason: '' })
const addProject = () => resumeContent.project.push({ name: '', role: '', start: '', end: '', description: '', tech_stack: [], url: '' })
const addSkill = () => resumeContent.skills.push({ category: '', name: '', proficiency: '' })
const updateTechStack = (item: ResumeProject, event: Event) => { item.tech_stack = (event.target as HTMLInputElement).value.split(/[,，]/).map((value) => value.trim()).filter(Boolean) }
const selectVersion = (version: ResumeVersion) => { resumeStore.setCurrentVersion(version); versionDrawerOpen.value = false }

const handleNameSave = async () => {
  if (demoMode.value || !resumeStore.currentResume || !editableName.value.trim()) return
  await resumeStore.update(resumeId.value, { name: editableName.value.trim() })
}
const handleSaveVersion = async () => {
  if (!newVersionNote.value.trim()) return message.warning('请输入版本备注')
  if (demoMode.value) { showSaveVersionModal.value = false; newVersionNote.value = ''; return message.success('本地预览：版本已模拟保存') }
  savingVersion.value = true
  const version = await resumeStore.createVersion(resumeId.value, { content: JSON.stringify(resumeContent), change_note: newVersionNote.value.trim() })
  savingVersion.value = false
  if (version) { showSaveVersionModal.value = false; newVersionNote.value = '' }
}
const previewNavigate = (target: 'interview' | 'delivery') => {
  message.info(target === 'interview' ? '面试训练场将在下一阶段完善' : '投递看板将在下一阶段完善')
}
const exportResume = () => window.print()
const backToList = () => demoMode.value ? router.push('/') : router.push('/app/resumes')
const formatVersionDate = (value: string) => value ? new Date(value).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' }) : ''

// 从个人资料 store 同步基础信息到简历 personal 字段
// 只在用户已填写个人资料时同步，且会覆盖示例默认值
const syncFromUserProfile = () => {
  if (!userProfileStore.hasFilled) return
  const b = userProfileStore.basic
  if (b.name) resumeContent.personal.name = b.name
  if (b.phone) resumeContent.personal.phone = b.phone
  if (b.email) resumeContent.personal.email = b.email
  if (b.city) resumeContent.personal.city = b.city
  if (b.github) resumeContent.personal.github = b.github
}

onMounted(async () => {
  document.documentElement.classList.add('resume-editor-scroll-lock')
  document.body.classList.add('resume-editor-scroll-lock')

  if (demoMode.value) return
  if (!resumeId.value || Number.isNaN(resumeId.value)) return router.push('/app/resumes')
  await resumeStore.fetchOne(resumeId.value)
  editableName.value = resumeStore.currentResume?.name || '未命名简历'
  targetJd.value = resumeStore.currentResume?.target_jd || targetJd.value
  await resumeStore.fetchVersions(resumeId.value)
  if (resumeStore.versions.length) resumeStore.setCurrentVersion(resumeStore.versions[0])
  // 加载完简历版本后，同步个人资料
  syncFromUserProfile()
})

onBeforeUnmount(() => {
  document.documentElement.classList.remove('resume-editor-scroll-lock')
  document.body.classList.remove('resume-editor-scroll-lock')
})

// 监听个人资料变化，实时同步到简历（用户在弹窗保存后立即生效）
watch(
  () => userProfileStore.basic,
  () => syncFromUserProfile(),
  { deep: true }
)
</script>

<style scoped>
:global(html.resume-editor-scroll-lock),
:global(body.resume-editor-scroll-lock) {
  height: 100%;
  overflow: hidden;
  overscroll-behavior: none;
}

/* ===== 局部变量映射到 Pinguo 全局 token（保持编辑器结构与全局视觉一致） ===== */
.resume-shell {
  --ink: var(--foreground);
  --muted: var(--muted-foreground);
  --line: var(--border);
  --green: var(--primary);
  --green-soft: var(--brand-50);
  width: 100%; height: 100%; min-height: 0; display: flex; overflow: hidden; background: var(--background-100); color: var(--ink);
  font-family: var(--font-sans);
}
.resume-shell.has-preview-sidebar { height: 100dvh; }
.product-sidebar { flex: 0 0 216px; min-width: 216px; display: flex; flex-direction: column; padding: 18px 12px 14px; border-right: 1px solid var(--sidebar-border); background: var(--sidebar); position: relative; z-index: 8; }
.sidebar-brand { width: 100%; display: flex; align-items: center; gap: 11px; padding: 4px 7px 20px; border: 0; background: transparent; color: var(--ink); text-align: left; cursor: pointer; }
.sidebar-logo { width: 36px; height: 36px; flex: 0 0 36px; display: grid; place-items: center; border-radius: 10px; background: var(--primary); color: var(--primary-foreground); font-size: 17px; font-weight: 850; }
.sidebar-brand > span:last-child, .product-nav button > span { min-width: 0; display: flex; flex-direction: column; }
.sidebar-brand strong { font-size: 16px; line-height: 1.2; }
.sidebar-brand small { margin-top: 3px; color: var(--muted-foreground); font-size: 10px; }
.product-nav { display: flex; flex-direction: column; gap: 6px; }
.product-nav button { width: 100%; min-height: 60px; display: flex; align-items: center; gap: 11px; padding: 10px 11px; border: 1px solid transparent; border-radius: 10px; background: transparent; color: var(--secondary-foreground); text-align: left; cursor: pointer; transition: background .18s, border-color .18s, color .18s; }
.product-nav button:hover { background: var(--sidebar-accent); color: var(--ink); }
.product-nav button.active { border-color: transparent; background: var(--card); color: var(--primary); box-shadow: var(--shadow-sm); }
.product-nav button > :deep(.anticon) { flex: 0 0 18px; font-size: 17px; }
.product-nav strong { font-size: 13px; line-height: 1.25; }
.product-nav small { margin-top: 4px; overflow: hidden; color: var(--muted-foreground); font-size: 9.5px; text-overflow: ellipsis; white-space: nowrap; }
.product-nav button.active small { color: var(--primary); }
.sidebar-foot { margin-top: auto; display: flex; align-items: center; gap: 7px; padding: 10px; color: var(--muted-foreground); font-size: 10px; }
.status-dot { width: 7px; height: 7px; border-radius: 50%; background: var(--success); box-shadow: 0 0 0 3px var(--state-success-surface); }
.resume-lab { height: 100%; min-height: 0; flex: 1 1 auto; min-width: 0; overflow: hidden; background: var(--background-100); color: var(--ink); }
.lab-toolbar { height: 64px; display: flex; align-items: center; justify-content: space-between; gap: 20px; padding: 0 24px; background: var(--card); border-bottom: 1px solid var(--line); position: relative; z-index: 5; }
.toolbar-leading, .toolbar-controls { display: flex; align-items: center; gap: 10px; min-width: 0; }
.toolbar-controls { overflow-x: auto; scrollbar-width: none; }
.icon-button { width: 36px; height: 36px; border: 0; background: transparent; border-radius: 9px; cursor: pointer; font-size: 16px; color: var(--muted-foreground); transition: background .18s, color .18s; }
.icon-button:hover { background: var(--background-200); color: var(--ink); }
.brand-mark { width: 34px; height: 34px; border-radius: 50%; display: grid; place-items: center; color: var(--primary-foreground); background: var(--primary); font-weight: 800; }
.resume-title-wrap { display: flex; align-items: center; gap: 12px; min-width: 0; }
.resume-name { width: 190px; border: 0; border-bottom: 1px solid transparent; outline: 0; color: var(--ink); font-size: 16px; font-weight: 750; background: transparent; }
.resume-name:focus { border-color: var(--primary); }
.save-state { color: var(--muted); font-size: 12px; white-space: nowrap; }
.save-state :deep(.anticon) { color: var(--primary); margin-right: 4px; }
.toolbar-select { height: 34px; display: flex; align-items: center; gap: 6px; padding: 0 9px; border-left: 1px solid var(--line); color: var(--secondary-foreground); }
.toolbar-select select { border: 0; outline: 0; background: transparent; font-size: 13px; cursor: pointer; color: var(--ink); }
.toolbar-select.compact { padding-right: 4px; }
.zoom-control { height: 34px; display: flex; align-items: center; border: 1px solid var(--line); border-radius: 9999px; overflow: hidden; background: var(--background-100); }
.zoom-control button { width: 30px; height: 100%; border: 0; background: transparent; cursor: pointer; color: var(--ink); }
.zoom-control span { width: 52px; text-align: center; font-size: 12px; color: var(--ink); }
.export-button { background: var(--primary) !important; border-color: var(--primary) !important; box-shadow: none; border-radius: 9999px !important; }
.export-button:hover { background: var(--brand-600) !important; border-color: var(--brand-600) !important; }
.lab-workspace { height: calc(100% - 64px); min-height: 0; display: grid; grid-template-columns: clamp(320px, 26vw, 400px) minmax(0, 1fr); gap: 20px; padding: 20px 24px 24px; background: var(--background-100); box-sizing: border-box; overflow: hidden; }
.editor-panel { background: var(--card); min-width: 0; overflow: hidden; border: 1px solid var(--line); border-radius: 16px; box-shadow: var(--shadow-sm); }
.smart-pane-wrap { display: flex; flex-direction: column; padding: 14px; overflow-y: auto; }
.panel-tabs { height: 51px; display: grid; grid-template-columns: 1fr 1fr; border-bottom: 1px solid var(--line); }
.panel-tabs button { border: 0; background: var(--card); color: var(--muted); font-weight: 650; cursor: pointer; position: relative; transition: color .18s; }
.panel-tabs button:hover { color: var(--ink); }
.panel-tabs button.active { color: var(--primary); }
.panel-tabs button.active::after { content: ''; height: 2px; background: var(--primary); position: absolute; left: 22px; right: 22px; bottom: 0; border-radius: 2px; }
.panel-hint { padding: 12px 16px; color: var(--muted); font-size: 12px; background: var(--background-100); border-bottom: 1px solid var(--line); }
.editor-scroll { height: calc(100% - 51px); overflow-y: auto; }
.form-section { border-bottom: 1px solid var(--line); }
.section-heading { width: 100%; height: 48px; padding: 0 9px 0 16px; display: flex; align-items: center; justify-content: space-between; background: var(--card); color: var(--ink); transition: opacity .18s, background .18s; }
.section-heading.disabled { opacity: .55; background: var(--background-100); }
.visibility-toggle { min-width: 0; flex: 1; height: 100%; display: flex; align-items: center; gap: 9px; cursor: pointer; font-weight: 720; }
.visibility-toggle input { position: absolute; width: 1px; height: 1px; opacity: 0; pointer-events: none; }
.custom-check { width: 18px; height: 18px; flex: 0 0 18px; display: grid; place-items: center; border: 1.5px solid var(--background-400); border-radius: 5px; background: var(--card); color: transparent; transition: background .18s, border-color .18s, box-shadow .18s; }
.custom-check :deep(.anticon) { font-size: 11px; }
.visibility-toggle input:checked + .custom-check { border-color: var(--primary); background: var(--primary); color: var(--primary-foreground); }
.visibility-toggle input:focus-visible + .custom-check { box-shadow: 0 0 0 3px rgba(0,122,255,.15); }
.visibility-toggle > :deep(.anticon) { color: var(--secondary-foreground); }
.visibility-toggle strong { overflow: hidden; font-size: 13px; text-overflow: ellipsis; white-space: nowrap; }
.collapse-button { width: 34px; height: 34px; flex: 0 0 34px; display: grid; place-items: center; border: 0; border-radius: 8px; background: transparent; color: var(--muted-foreground); cursor: pointer; transition: background .18s, color .18s; }
.collapse-button:hover { background: var(--background-200); color: var(--ink); }
.collapse-button :deep(.anticon) { transition: transform .2s; }
.form-section.collapsed .collapse-button :deep(.anticon) { transform: rotate(-90deg); }
.form-section.collapsed .section-body { display: none; }
.section-body { padding: 0 16px 16px; }
.grid-two { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.field { display: flex; flex-direction: column; gap: 5px; min-width: 0; }
.field.full { grid-column: 1 / -1; }
.field span { color: var(--muted); font-size: 11px; font-weight: 650; }
.field input, .field textarea, .skill-row input { width: 100%; border: 1px solid var(--line); border-radius: 8px; outline: none; padding: 8px 9px; color: var(--ink); background: var(--card); font: inherit; font-size: 12px; transition: border-color .2s, box-shadow .2s; }
.field textarea { resize: vertical; line-height: 1.55; }
.field input:focus, .field textarea:focus, .skill-row input:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(0,122,255,.12); }
.field input::placeholder, .field textarea::placeholder { color: var(--muted-foreground); }
.repeat-list { display: flex; flex-direction: column; gap: 10px; }
.repeat-item { position: relative; padding: 12px; border: 1px solid var(--line); border-radius: 12px; background: var(--background-100); }
.remove-item { position: absolute; top: 5px; right: 5px; width: 26px; height: 26px; border: 0; border-radius: 6px; background: transparent; color: var(--muted-foreground); cursor: pointer; z-index: 2; transition: background .15s, color .15s; }
.remove-item:hover { color: var(--destructive); background: var(--state-error-surface); }
.add-item { min-height: 34px; border: 1px dashed var(--background-400); border-radius: 9999px; background: var(--card); color: var(--primary); cursor: pointer; font-weight: 650; transition: background .15s, border-color .15s; }
.add-item:hover { background: var(--brand-50); border-color: var(--primary); }
.skill-row { display: grid; grid-template-columns: 90px 1fr 26px; gap: 6px; }
.skill-row button { border: 0; background: transparent; color: var(--muted-foreground); cursor: pointer; }
.skill-row button:hover { color: var(--destructive); }
.preview-stage { min-width: 0; overflow: hidden; background: var(--card); display: flex; flex-direction: column; border: 1px solid var(--line); border-radius: 16px; box-shadow: var(--shadow-sm); }
.stage-meta { height: 42px; padding: 0 18px; display: flex; align-items: center; justify-content: space-between; color: var(--muted-foreground); font-size: 11px; border-bottom: 1px solid var(--line); }
.paper-viewport { flex: 1; overflow: auto; padding: 26px 32px 80px; display: flex; justify-content: center; align-items: flex-start; background: var(--background-100); }
.resume-paper { width: 794px; min-height: 1123px; transform-origin: top center; padding: 56px 60px; background: var(--card); box-shadow: var(--shadow-lg); color: var(--text-800); transition: transform .18s ease; border-radius: 4px; }
.resume-paper.font-serif { font-family: 'Songti SC', 'Noto Serif SC', serif; }
.resume-paper.font-sans { font-family: 'PingFang SC', 'Noto Sans SC', sans-serif; }
.paper-header { padding-bottom: 16px; border-bottom: 2px solid #26312b; }
.paper-header h1 { margin: 0; font-size: 31px; line-height: 1.15; letter-spacing: .08em; font-weight: 800; }
.paper-role { margin: 8px 0 10px; color: var(--green); font-size: 14px; font-weight: 750; }
.contact-line { display: flex; flex-wrap: wrap; gap: 5px 18px; color: #4d5651; font-size: 10.5px; }
.contact-line span:not(:last-child)::after { content: '·'; margin-left: 18px; color: #9ca39f; }
.paper-section { margin-top: 20px; }
.paper-section h2 { margin: 0 0 9px; padding-bottom: 5px; border-bottom: 1px solid #39433e; font-size: 15px; line-height: 1; letter-spacing: .08em; }
.paper-section p { margin: 0; font-size: 10.8px; line-height: 1.65; }
.resume-entry + .resume-entry { margin-top: 13px; }
.entry-heading, .entry-subheading { display: flex; align-items: baseline; justify-content: space-between; gap: 12px; }
.entry-heading strong { font-size: 11.5px; }
.entry-heading span, .entry-subheading { color: #545e58; font-size: 10px; }
.entry-subheading { margin-top: 3px; justify-content: flex-start; }
.entry-subheading .tech-stack { margin-left: auto; color: #68736d; }
.entry-note { color: #505a54; }
.resume-entry ul { margin: 6px 0 0; padding-left: 17px; }
.resume-entry li { margin: 2px 0; font-size: 10.4px; line-height: 1.55; }
.skill-section p + p { margin-top: 3px; }
.template-modern .paper-header { padding: 18px 20px; border: 0; background: #edf5f0; }
.template-modern .paper-section h2 { color: var(--green); border-color: #8fb9a3; }
.template-executive { color: #202a35; font-family: Georgia, 'Songti SC', serif !important; }
.template-executive .paper-header { text-align: center; border-bottom: 1px solid #28384d; }
.template-executive .paper-header h1 { color: #28384d; letter-spacing: .14em; }
.template-executive .paper-role { color: #667689; }
.template-executive .contact-line { justify-content: center; }
.template-executive .paper-section h2 { border: 0; color: #28384d; text-align: center; letter-spacing: .16em; }
.template-executive .paper-section h2::after { content: ''; width: 34px; height: 1px; display: block; margin: 7px auto 0; background: #9ba6b2; }
.template-compact { padding: 42px 48px; }
.template-compact .paper-header { padding-bottom: 10px; }
.template-compact .paper-header h1 { font-size: 27px; }
.template-compact .paper-section { margin-top: 13px; }
.template-compact .paper-section h2 { margin-bottom: 6px; padding-bottom: 4px; font-size: 13px; }
.template-compact .resume-entry + .resume-entry { margin-top: 7px; }
.template-compact .resume-entry li { margin: 1px 0; line-height: 1.42; }
.template-sidebar {
  display: grid;
  grid-template-columns: 176px minmax(0, 1fr);
  column-gap: 34px;
  align-content: start;
  padding: 0 48px 52px 0;
  background: linear-gradient(90deg, #214c3d 0 176px, #fff 176px);
}
.template-sidebar .paper-header {
  grid-column: 1;
  grid-row: 1 / span 20;
  min-height: 1123px;
  padding: 56px 23px;
  border: 0;
  color: #fff;
}
.template-sidebar .paper-header h1 { font-size: 25px; letter-spacing: .03em; }
.template-sidebar .paper-role { color: #bfe1d2; }
.template-sidebar .contact-line { flex-direction: column; gap: 8px; color: #e7f2ed; line-height: 1.45; }
.template-sidebar .contact-line span::after { display: none; }
.template-sidebar .paper-section { grid-column: 2; margin-top: 20px; }
.template-sidebar .paper-section:first-of-type { margin-top: 56px; }
.template-sidebar .paper-section h2 { color: #214c3d; border-color: #8eb2a2; }
.template-editorial .paper-header { position: relative; padding-bottom: 24px; border: 0; }
.template-editorial .paper-header::after { content: ''; width: 72px; height: 5px; position: absolute; bottom: 0; left: 0; background: #b5523b; }
.template-editorial .paper-header h1 { max-width: 520px; color: #9d4532; font-family: Georgia, 'Songti SC', serif; font-size: 44px; line-height: .95; letter-spacing: -.02em; }
.template-editorial .paper-role { color: #3f4944; letter-spacing: .08em; }
.template-editorial .paper-section h2 { border: 0; color: #a04935; font-family: Georgia, 'Songti SC', serif; font-size: 18px; }
.template-minimal { padding: 78px 74px; }
.template-minimal .paper-header { padding-bottom: 28px; border-bottom: 1px solid #a9b0ac; }
.template-minimal .paper-header h1 { font-size: 29px; font-weight: 550; letter-spacing: .18em; }
.template-minimal .paper-role { color: #555f5a; font-weight: 550; }
.template-minimal .paper-section { margin-top: 29px; }
.template-minimal .paper-section h2 { border: 0; font-size: 12px; font-weight: 650; letter-spacing: .22em; }
.template-academic { padding: 50px 54px; font-family: Georgia, 'Songti SC', serif !important; color: #262522; }
.template-academic .paper-header { text-align: center; border-bottom: 2px double #4e4b46; }
.template-academic .paper-header h1 { font-size: 28px; font-weight: 600; }
.template-academic .paper-role { color: #4d4a45; font-style: italic; font-weight: 500; }
.template-academic .contact-line { justify-content: center; }
.template-academic .paper-section h2 { border-bottom: 1px solid #59554f; font-size: 13px; text-transform: uppercase; }
.template-creative { padding-top: 0; }
.template-creative .paper-header { margin: 0 -60px 28px; padding: 48px 60px 34px; border: 0; background: #126b51; color: #fff; }
.template-creative .paper-header h1 { font-size: 39px; letter-spacing: -.015em; }
.template-creative .paper-role { color: #c9eadc; }
.template-creative .contact-line { color: #edf8f3; }
.template-creative .paper-section h2 { padding: 7px 10px; border: 0; background: #e7f2ed; color: #126b51; }
.template-graduate { position: relative; padding-left: 72px; }
.template-graduate::before { content: ''; width: 8px; position: absolute; inset: 0 auto 0 0; background: #277c87; }
.template-graduate .paper-header { border-color: #277c87; }
.template-graduate .paper-header h1 { color: #205f68; }
.template-graduate .paper-role { color: #277c87; }
.template-graduate .paper-section h2 { color: #277c87; border-color: #9bc1c6; }
.density-compact { padding-top: 46px; padding-bottom: 46px; }
.density-compact .paper-section { margin-top: 15px; }
.density-compact .resume-entry + .resume-entry { margin-top: 9px; }
.ai-heading { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 16px; }
.ai-heading h2 { margin: 3px 0 0; font-size: 15px; }
.ai-heading > :deep(.anticon) { color: var(--green); font-size: 20px; }
.ai-kicker { color: var(--green); font-size: 12px; font-weight: 800; letter-spacing: .08em; }

/* === 智能完善 tabs === */
.smart-tabs { display: grid; grid-template-columns: repeat(3, 1fr); gap: 5px; padding: 5px; background: var(--background-200); border-radius: 12px; margin-bottom: 16px; }
.smart-tab-btn { width: 100%; height: 42px; min-width: 0; display: grid; grid-template-columns: 18px auto; align-items: center; justify-content: center; column-gap: 6px; padding: 0 8px; border: 0; background: transparent; color: var(--muted); font-size: 11px; line-height: 18px; font-weight: 650; cursor: pointer; border-radius: 9px; transition: background .15s, color .15s, box-shadow .15s; }
.smart-tab-btn:hover:not(.active) { background: rgba(255, 255, 255, 0.55); color: var(--secondary-foreground); }
.smart-tab-btn.active { background: var(--card); color: var(--primary); box-shadow: var(--shadow-sm); }
.smart-tab-btn :deep(.anticon) { width: 18px; height: 18px; display: grid; place-items: center; font-size: 15px; line-height: 1; }
.smart-tab-btn > span:last-child { min-width: 0; line-height: 18px; white-space: nowrap; }

.smart-pane { display: flex; flex-direction: column; gap: 14px; }
.pane-head { min-height: 30px; display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 0 2px; }
.pane-title { font-size: 13px; font-weight: 700; color: var(--ink); }
.mini-btn { display: inline-flex; align-items: center; gap: 4px; padding: 5px 10px; border: 1px solid var(--line); border-radius: 9999px; background: var(--card); color: var(--muted); font-size: 11px; cursor: pointer; transition: all .15s; }
.mini-btn:hover:not(:disabled) { border-color: var(--primary); color: var(--primary); }
.mini-btn.primary { background: var(--primary); color: var(--primary-foreground); border-color: var(--primary); }
.mini-btn.primary:hover:not(:disabled) { background: var(--brand-600); border-color: var(--brand-600); }
.mini-btn:disabled { opacity: 0.6; cursor: not-allowed; }

/* === 分析页 === */
.score-overview { display: flex; align-items: center; gap: 14px; padding: 14px; background: linear-gradient(135deg, var(--brand-50) 0%, var(--card) 100%); border: 1px solid var(--line); border-radius: 12px; }
.score-ring { position: relative; width: 80px; height: 80px; flex-shrink: 0; }
.score-ring svg { width: 100%; height: 100%; }
.score-num { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; color: var(--primary); font-size: 22px; font-weight: 800; }
.score-num small { font-size: 11px; font-weight: 600; color: var(--muted); margin-left: 2px; }
.score-meta { flex: 1; min-width: 0; }
.score-meta strong { display: block; font-size: 14px; color: var(--ink); margin-bottom: 4px; }
.score-meta p { margin: 0; font-size: 11px; color: var(--muted); line-height: 1.5; }

.dim-card { padding: 12px 14px; border: 1px solid var(--line); border-radius: 12px; background: var(--card); }
.dim-card h4 { margin: 0 0 10px; font-size: 12px; font-weight: 700; color: var(--ink); }
.dim-card h4 small { font-weight: 500; color: var(--muted); margin-left: 6px; }
.dim-row { display: grid; grid-template-columns: 80px 1fr 28px; align-items: center; gap: 8px; margin: 7px 0; font-size: 11px; }
.dim-label { color: var(--muted); }
.dim-track { height: 5px; border-radius: 9999px; background: var(--background-200); overflow: hidden; }
.dim-track i { display: block; height: 100%; border-radius: inherit; transition: width .4s ease; }
.dim-row b { font-size: 11px; text-align: right; }

.hl-list { display: flex; flex-direction: column; gap: 8px; }
.hl-item { display: flex; gap: 8px; padding: 9px 0; border-top: 1px solid var(--background-200); }
.hl-item:first-child { border-top: 0; padding-top: 0; }
.hl-item :deep(.anticon) { flex-shrink: 0; font-size: 14px; margin-top: 1px; }
.hl-good :deep(.anticon) { color: var(--success); }
.hl-bad :deep(.anticon) { color: var(--warning); }
.hl-item strong { font-size: 12px; color: var(--ink); display: block; margin-bottom: 3px; }
.hl-item p { margin: 0; font-size: 10.5px; color: var(--muted); line-height: 1.5; }
.link-btn { display: inline-flex; align-items: center; gap: 3px; margin-top: 6px; padding: 4px 9px; border: 0; background: var(--brand-50); color: var(--primary); font-size: 10px; font-weight: 700; border-radius: 9999px; cursor: pointer; transition: background .15s, color .15s; }
.link-btn:hover { background: var(--primary); color: var(--primary-foreground); }
.link-btn:disabled { opacity: 0.5; cursor: not-allowed; }

/* === JD 匹配页 === */
.jd-textarea-wrap { position: relative; display: block; }
.jd-textarea-wrap textarea { width: 100%; resize: vertical; font-size: 11.5px; line-height: 1.55; padding: 10px 11px; border: 1px solid var(--line); border-radius: 8px; outline: none; font: inherit; color: var(--ink); background: var(--card); min-height: 110px; transition: border-color .2s, box-shadow .2s; }
.jd-textarea-wrap textarea:focus { border-color: var(--primary); box-shadow: 0 0 0 3px rgba(0,122,255,.12); }
.jd-textarea-wrap textarea::placeholder { color: var(--muted-foreground); }
.jd-counter { position: absolute; right: 9px; bottom: 7px; font-size: 10px; color: var(--muted); background: var(--card); padding: 1px 5px; border-radius: 4px; }

.jd-result { display: flex; flex-direction: column; gap: 14px; padding: 14px; background: var(--background-100); border: 1px solid var(--line); border-radius: 12px; }
.jd-match-ring { position: relative; width: 120px; height: 120px; margin: 0 auto; }
.jd-match-ring svg { width: 100%; height: 100%; }
.jd-match-num { position: absolute; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center; }
.jd-match-num strong { font-size: 24px; font-weight: 800; color: var(--ink); }
.jd-match-num small { font-size: 10px; color: var(--muted); margin-top: 2px; }

.jd-stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.jd-stat-item { display: flex; flex-direction: column; align-items: center; padding: 8px; background: var(--card); border-radius: 10px; }
.jd-stat-num { font-size: 18px; font-weight: 800; }
.jd-stat-num.green { color: var(--success); }
.jd-stat-num.red { color: var(--destructive); }
.jd-stat-num.blue { color: var(--primary); }
.jd-stat-label { font-size: 10px; color: var(--muted); margin-top: 2px; }

.kw-section h5 { margin: 0 0 8px; font-size: 11px; font-weight: 700; color: var(--ink); }
.kw-section h5 small { font-weight: 500; color: var(--muted); margin-left: 4px; }
.kw-tags { display: flex; flex-wrap: wrap; gap: 5px; }
.kw-tag { padding: 3px 9px; border-radius: 9999px; font-size: 10.5px; font-weight: 600; }
.kw-green { background: rgba(52,199,89,0.12); color: var(--success); }
.kw-red { background: var(--state-error-surface); color: var(--destructive); }
.kw-empty { color: var(--muted); font-size: 11px; }

.suggest-item { display: flex; gap: 8px; padding: 9px 0; border-top: 1px solid var(--background-200); }
.suggest-item:first-child { border-top: 0; }
.suggest-item b { flex: 0 0 20px; height: 20px; display: grid; place-items: center; border-radius: 50%; background: var(--brand-50); color: var(--primary); font-size: 10px; }
.suggest-item strong { font-size: 11.5px; color: var(--ink); display: block; }
.suggest-item p { margin: 3px 0 0; font-size: 10px; color: var(--muted); line-height: 1.5; }

.jd-empty { min-height: 82px; display: flex; align-items: center; justify-content: center; gap: 10px; padding: 18px 20px; border: 1px dashed var(--border); border-radius: 12px; background: var(--background-100); color: var(--muted); }
.jd-empty :deep(.anticon) { flex: 0 0 auto; font-size: 24px; opacity: 0.55; }
.jd-empty p { margin: 0; font-size: 11px; line-height: 1.5; }

/* === AI 优化页 === */
.opt-section { padding: 12px 14px; border: 1px solid var(--line); border-radius: 12px; background: var(--card); margin-bottom: 10px; }
.opt-section h4 { margin: 0 0 6px; font-size: 12px; font-weight: 700; }
.opt-desc { margin: 0 0 10px; font-size: 10.5px; color: var(--muted); line-height: 1.5; }
.opt-actions { display: flex; gap: 8px; }
.opt-btn { flex: 1; padding: 9px 12px; border: 0; border-radius: 9999px; background: var(--primary); color: var(--primary-foreground); font-size: 12px; font-weight: 700; cursor: pointer; transition: background .15s; }
.opt-btn:hover:not(:disabled) { background: var(--brand-600); }
.opt-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.opt-item { padding: 10px 0; border-top: 1px solid var(--background-200); }
.opt-item:first-of-type { border-top: 0; padding-top: 0; }
.opt-item-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 4px; }
.opt-item-head strong { font-size: 12px; color: var(--ink); }
.opt-item-status { padding: 2px 7px; border-radius: 9999px; font-size: 10px; font-weight: 600; }
.st-idle { background: var(--background-200); color: var(--muted); }
.st-done { background: var(--state-success-surface); color: var(--success); }
.opt-item-desc { margin: 0 0 6px; font-size: 10.5px; color: var(--muted); line-height: 1.5; }

.privacy-note { position: sticky; bottom: -14px; z-index: 2; display: flex; align-items: center; justify-content: center; gap: 5px; margin: auto -14px -14px; padding: 12px 14px 14px; border-top: 10px solid transparent; background: var(--card); background-clip: padding-box; box-shadow: inset 0 1px 0 var(--background-200); color: var(--muted-foreground); font-size: 9.5px; line-height: 1.45; text-align: center; }
.version-list { display: flex; flex-direction: column; gap: 8px; }
.version-list button { display: flex; flex-direction: column; gap: 5px; padding: 12px; border: 1px solid var(--line); border-radius: 10px; background: var(--card); text-align: left; cursor: pointer; transition: border-color .15s, background .15s; }
.version-list button:hover { background: var(--background-100); }
.version-list button.active { border-color: var(--primary); background: var(--brand-50); }
.version-list button span { display: flex; justify-content: space-between; }
.version-list small, .version-list em { color: var(--muted); font-size: 11px; font-style: normal; }
@media (max-width: 1380px) { .lab-workspace { grid-template-columns: clamp(300px, 30vw, 360px) minmax(0, 1fr); } .toolbar-select { display: none; } .save-state { display: none; } }
@media (max-width: 1060px) { .product-sidebar { flex-basis: 176px; min-width: 176px; padding-inline: 9px; } .product-nav small, .sidebar-foot { display: none; } .product-nav button { min-height: 50px; } .lab-toolbar { padding-inline: 14px; gap: 10px; } .resume-name { width: 150px; } }
@media (max-width: 820px) {
  .resume-shell { height: 100%; min-height: 0; flex-direction: column; overflow: hidden; }
  .resume-shell.has-preview-sidebar { height: 100dvh; }
  .product-sidebar { width: 100%; min-width: 0; flex: 0 0 58px; min-height: 58px; flex-direction: row; align-items: center; padding: 7px 10px; border-right: 0; border-bottom: 1px solid var(--sidebar-border); }
  .sidebar-brand { width: auto; padding: 0 8px 0 0; }
  .sidebar-logo { width: 34px; height: 34px; flex-basis: 34px; }
  .product-nav { flex: 1; min-width: 0; flex-direction: row; justify-content: flex-end; gap: 3px; }
  .product-nav button { width: auto; min-width: 44px; min-height: 42px; flex: 1 1 0; justify-content: center; padding: 6px 8px; }
  .product-nav button > span { display: flex; }
  .product-nav button small { display: none; }
  .product-nav button > :deep(.anticon) { font-size: 15px; }
  .product-nav strong { font-size: 11px; white-space: nowrap; }
  .resume-lab { height: calc(100dvh - 58px); min-height: 0; }
  .lab-toolbar { height: 108px; min-height: 108px; padding: 9px 12px; flex-wrap: wrap; align-content: center; }
  .toolbar-leading { width: 100%; }
  .toolbar-controls { width: 100%; overflow-x: auto; padding-bottom: 2px; }
  .resume-name { width: min(46vw, 190px); }
  .lab-workspace { height: calc(100% - 108px); display: grid; grid-template-columns: minmax(270px, 38vw) minmax(0, 1fr); gap: 12px; padding: 12px 14px 14px; }
  .editor-panel { height: auto; border-right: 1px solid var(--line); border-bottom: 0; }
  .preview-stage { height: auto; min-height: 0; }
  .paper-viewport { padding: 18px 18px 60px; justify-content: flex-start; }
  .resume-paper { flex: 0 0 794px; }
}
@media (max-width: 620px) {
  .product-nav button { padding-inline: 5px; }
  .product-nav button > :deep(.anticon) { display: none; }
  .lab-workspace { grid-template-columns: 1fr; grid-template-rows: minmax(270px, 46%) minmax(0, 54%); gap: 10px; padding: 10px 12px 12px; }
  .editor-panel { border-right: 0; border-bottom: 1px solid var(--line); }
  .stage-meta span:last-child { display: none; }
}
@media (max-height: 680px) and (min-width: 821px) { .lab-toolbar { height: 54px; padding-inline: 14px; } .lab-workspace { height: calc(100% - 54px); } .panel-tabs { height: 44px; } .editor-scroll { height: calc(100% - 44px); } }
@media print { .product-sidebar, .lab-toolbar, .editor-panel, .stage-meta { display: none !important; } .resume-shell, .resume-lab, .lab-workspace, .preview-stage, .paper-viewport { display: block; height: auto; overflow: visible; background: #fff; padding: 0; } .resume-paper { transform: none !important; box-shadow: none; margin: 0; width: 210mm; min-height: 297mm; } }
</style>
