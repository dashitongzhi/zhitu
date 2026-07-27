<template>
  <div class="dk-shell" tabindex="0" @keydown="handleKeydown">
    <!-- ==================== Quick Capture 顶部栏 ==================== -->
    <div class="qc-bar">
      <div class="qc-input-wrap">
        <ThunderboltOutlined class="qc-icon" />
        <input
          v-model="quickCapture"
          class="qc-input"
          placeholder="粘贴 JD 快速创建投递，或输入「公司 / 岗位」后回车…"
          @keydown.enter="handleQuickCapture"
          @paste="handlePaste"
        />
        <kbd class="qc-kbd">⏎</kbd>
      </div>
      <div class="qc-hint" v-if="quickCaptureParsed.company || quickCaptureParsed.position">
        <span class="qc-parsed-tag" v-if="quickCaptureParsed.company">
          <BankOutlined /> {{ quickCaptureParsed.company }}
        </span>
        <span class="qc-parsed-tag" v-if="quickCaptureParsed.position">
          <SolutionOutlined /> {{ quickCaptureParsed.position }}
        </span>
      </div>
    </div>

    <!-- ==================== 统计卡片网格 ==================== -->
    <div class="stat-section-title">个人数据</div>
    <div class="stat-grid">
      <div
        v-for="card in statCards"
        :key="card.key"
        class="stat-card"
        :class="card.accentClass"
      >
        <div class="stat-icon-wrap" :style="{ background: card.iconBg, color: card.iconColor }">
          <component :is="card.icon" />
        </div>
        <div class="stat-body">
          <div class="stat-number">{{ card.value }}</div>
          <div class="stat-label">{{ card.label }}</div>
          <div class="stat-change" :class="card.changeClass">{{ card.change }}</div>
        </div>
        <svg class="stat-spark" viewBox="0 0 80 24" preserveAspectRatio="none">
          <polyline
            :points="card.sparkline"
            fill="none"
            :stroke="card.sparkColor"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </div>
    </div>

    <!-- ==================== 视图切换 + 筛选器 ==================== -->
    <div class="section-row">
      <div class="view-toggle-group">
        <button
          class="view-toggle-btn"
          :data-active="viewMode === 'personal'"
          @click="viewMode = 'personal'"
        >
          <TableOutlined /> 个人投递
        </button>
        <button
          class="view-toggle-btn"
          :data-active="viewMode === 'platform'"
          @click="viewMode = 'platform'"
        >
          <PieChartOutlined /> 平台数据
        </button>
      </div>

      <div v-if="viewMode !== 'platform'" class="search-wrap">
        <SearchOutlined class="search-icon" />
        <input
          v-model="searchKeyword"
          class="search-input"
          placeholder="搜索公司、岗位…"
        />
      </div>

      <div v-if="viewMode !== 'platform'" class="filter-cluster">
        <select class="filter-select" v-model="filters.status" @change="handleFilter">
          <option value="">全部状态</option>
          <option v-for="s in statusList" :key="s.value" :value="s.value">{{ s.label }}</option>
        </select>
        <select class="filter-select" v-model="filters.priority" @change="handleFilter">
          <option value="">全部优先级</option>
          <option value="high">高</option>
          <option value="medium">中</option>
          <option value="low">低</option>
        </select>
        <select class="filter-select" v-model="filters.channel" @change="handleFilter">
          <option value="">全部渠道</option>
          <option v-for="c in channelOptions" :key="c.value" :value="c.value">{{ c.label }}</option>
        </select>
        <button class="btn-primary-capsule" @click="showCreateModal">
          <PlusOutlined /> 新增投递
        </button>
      </div>
    </div>

    <!-- ==================== 主体：数据区 + 详情面板 ==================== -->
    <div
      class="content-flex"
      :class="{ 'platform-mode': viewMode === 'platform' }"
      v-if="filteredDeliveries.length > 0 || hasActiveFilter || viewMode === 'platform'"
    >
      <div class="content-main">
        <!-- 空状态 -->
        <div v-if="filteredDeliveries.length === 0 && viewMode !== 'platform'" class="empty-state">
          <div class="empty-illustration">
            <InboxOutlined />
          </div>
          <div class="empty-title">没有匹配的投递记录</div>
          <div class="empty-desc">尝试调整筛选条件，或清空所有筛选查看全部记录</div>
          <button class="btn-secondary-capsule" @click="clearFilters">清空筛选</button>
        </div>

        <!-- 视图切换动画包装 -->
        <Transition name="view-fade" mode="out-in">
          <!-- ===== 个人投递视图（表格 + 移动端卡片） ===== -->
          <div v-if="viewMode === 'personal' && filteredDeliveries.length > 0" key="personal" class="view-wrap">
            <!-- 桌面/平板表格 -->
            <div class="kb-table-card kb-table-desktop">
              <div class="kb-table-scroll">
                <table class="kb-table">
                  <thead>
                    <tr>
                      <th>公司</th>
                      <th>岗位</th>
                      <th>渠道</th>
                      <th>投递日期</th>
                      <th>状态</th>
                      <th>面试进度</th>
                      <th>下次面试</th>
                      <th>HR 最新反馈</th>
                      <th>优先级</th>
                      <th style="text-align: right;">操作</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="(item, idx) in filteredDeliveries"
                      :key="item.id"
                      :class="{ selected: selectedId === item.id }"
                      :data-idx="idx"
                      @click="selectDelivery(item)"
                      @mouseenter="hoveredId = item.id"
                      @mouseleave="hoveredId = null"
                    >
                      <td>
                        <div class="kb-company-cell">
                          <span class="kb-company">{{ item.company }}</span>
                          <span v-if="item.remark" class="kb-remark-dot" :title="item.remark">·</span>
                        </div>
                      </td>
                      <td>{{ item.position }}</td>
                      <td><span class="kb-channel" :class="getChannelClass(item.channel)">{{ getChannelLabel(item.channel) }}</span></td>
                      <td class="mono">{{ formatDate(item.apply_date) }}</td>
                      <td>
                        <span class="kb-status" :class="getStatusClass(item.status)">
                          <span class="status-dot" :style="{ background: getStatusColor(item.status) }"></span>
                          {{ getStatusText(item.status) }}
                        </span>
                      </td>
                      <td>
                        <div class="kb-dots" :title="getProgressTooltip(item)">
                          <span
                            v-for="(d, i) in getProgressDots(item)"
                            :key="i"
                            class="dot"
                            :class="d"
                          ></span>
                        </div>
                      </td>
                      <td>
                        <span v-if="getNextInterview(item)" class="kb-next-interview">
                          <CalendarOutlined />
                          {{ formatShortDateTime(getNextInterview(item).interview_time) }}
                        </span>
                        <span v-else class="kb-dash">—</span>
                      </td>
                      <td class="kb-feedback-cell">
                        <span v-if="getLatestFeedback(item.id)">{{ truncate(getLatestFeedback(item.id).summary, 24) }}</span>
                        <span v-else class="kb-dash">—</span>
                      </td>
                      <td><span class="kb-priority" :class="getPriorityClass(item.priority)">{{ getPriorityText(item.priority) }}</span></td>
                      <td style="text-align: right; white-space: nowrap;" @click.stop>
                        <button class="kb-action" @click="selectDelivery(item)" title="查看详情">
                          <RightOutlined />
                        </button>
                        <a-dropdown :trigger="['click']">
                          <button class="kb-action" title="更多操作"><EllipsisOutlined /></button>
                          <template #overlay>
                            <a-menu>
                              <a-menu-item
                                v-for="s in getAvailableTransitions(item.status)"
                                :key="s.value"
                                @click="handleStatusChange(item, s.value)"
                              >
                                <RightCircleOutlined /> 转为{{ s.label }}
                              </a-menu-item>
                              <a-menu-divider />
                              <a-menu-item @click="showRoundModalFor(item)"><PlusOutlined /> 添加轮次</a-menu-item>
                              <a-menu-item @click="showFeedbackModalFor(item)"><MessageOutlined /> 添加反馈</a-menu-item>
                              <a-menu-divider />
                              <a-menu-item danger @click="handleDelete(item.id)"><DeleteOutlined /> 删除</a-menu-item>
                            </a-menu>
                          </template>
                        </a-dropdown>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- 移动端卡片堆叠 -->
            <div class="kb-cards-mobile">
              <div
                v-for="item in filteredDeliveries"
                :key="item.id"
                class="kb-mobile-card"
                :class="{ selected: selectedId === item.id }"
                @click="selectDelivery(item)"
              >
                <div class="kb-mobile-head">
                  <span class="kb-company">{{ item.company }}</span>
                  <span class="kb-status" :class="getStatusClass(item.status)">
                    <span class="status-dot" :style="{ background: getStatusColor(item.status) }"></span>
                    {{ getStatusText(item.status) }}
                  </span>
                </div>
                <div class="kb-mobile-pos">{{ item.position }}</div>
                <div class="kb-mobile-meta">
                  <span class="kb-channel" :class="getChannelClass(item.channel)">{{ getChannelLabel(item.channel) }}</span>
                  <span class="mono">{{ formatDate(item.apply_date) }}</span>
                  <span class="kb-priority" :class="getPriorityClass(item.priority)">{{ getPriorityText(item.priority) }}</span>
                </div>
                <div class="kb-mobile-foot" v-if="getNextInterview(item) || getLatestFeedback(item.id)">
                  <span v-if="getNextInterview(item)" class="kb-next-interview">
                    <CalendarOutlined /> {{ formatShortDateTime(getNextInterview(item).interview_time) }}
                  </span>
                  <span v-if="getLatestFeedback(item.id)" class="kb-mobile-fb">{{ getLatestFeedback(item.id).summary }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- ===== 平台数据视图 ===== -->
          <div v-else-if="viewMode === 'platform'" key="platform" class="platform-view">
            <!-- 顶部核心指标 -->
            <div class="platform-hero">
              <div class="platform-hero-card">
                <div class="platform-hero-icon" style="background: var(--brand-50); color: var(--primary);">
                  <TeamOutlined />
                </div>
                <div>
                  <div class="platform-hero-value">{{ platformStats.totalUsers.toLocaleString() }}</div>
                  <div class="platform-hero-label">累计服务用户</div>
                </div>
              </div>
              <div class="platform-hero-card">
                <div class="platform-hero-icon" style="background: var(--state-success-surface); color: var(--success);">
                  <TrophyOutlined />
                </div>
                <div>
                  <div class="platform-hero-value">{{ platformStats.successRate }}%</div>
                  <div class="platform-hero-label">综合成功率</div>
                </div>
              </div>
              <div class="platform-hero-card">
                <div class="platform-hero-icon" style="background: rgba(175,82,222,0.14); color: #8e3db5;">
                  <SendOutlined />
                </div>
                <div>
                  <div class="platform-hero-value">{{ platformStats.totalDeliveries.toLocaleString() }}</div>
                  <div class="platform-hero-label">累计投递次数</div>
                </div>
              </div>
              <div class="platform-hero-card">
                <div class="platform-hero-icon" style="background: rgba(255,149,0,0.14); color: #c46900;">
                  <RiseOutlined />
                </div>
                <div>
                  <div class="platform-hero-value">{{ platformStats.avgCycleDays }}</div>
                  <div class="platform-hero-label">平均求职周期（天）</div>
                </div>
              </div>
            </div>

            <div class="platform-grid-2">
              <!-- 饼图：投递结果分布 -->
              <div class="platform-card">
                <div class="platform-card-head">
                  <div class="platform-card-title">
                    <PieChartOutlined /> 投递结果分布
                  </div>
                  <span class="platform-card-sub">基于全平台累计数据</span>
                </div>
                <div class="pie-wrap">
                  <svg class="pie-svg" viewBox="0 0 200 200">
                    <circle
                      v-for="(seg, i) in pieSegments"
                      :key="i"
                      cx="100" cy="100" r="70"
                      :fill="seg.fill"
                      :stroke="seg.stroke"
                      :stroke-width="seg.strokeWidth"
                      :stroke-dasharray="seg.dashArray"
                      :stroke-dashoffset="seg.dashOffset"
                      transform="rotate(-90 100 100)"
                    />
                    <text x="100" y="92" text-anchor="middle" class="pie-center-num">{{ platformStats.totalDeliveries.toLocaleString() }}</text>
                    <text x="100" y="112" text-anchor="middle" class="pie-center-label">累计投递</text>
                  </svg>
                  <div class="pie-legend">
                    <div v-for="(seg, i) in pieSegments" :key="i" class="pie-legend-item">
                      <span class="pie-legend-dot" :style="{ background: seg.color }"></span>
                      <span class="pie-legend-name">{{ seg.label }}</span>
                      <span class="pie-legend-value">{{ seg.value.toLocaleString() }}</span>
                      <span class="pie-legend-pct">{{ seg.pct }}%</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 饼图：服务行业分布 -->
              <div class="platform-card">
                <div class="platform-card-head">
                  <div class="platform-card-title">
                    <SolutionOutlined /> 服务行业分布
                  </div>
                  <span class="platform-card-sub">用户所在行业 Top 6</span>
                </div>
                <div class="pie-wrap">
                  <svg class="pie-svg" viewBox="0 0 200 200">
                    <circle
                      v-for="(seg, i) in industryPieSegments"
                      :key="i"
                      cx="100" cy="100" r="70"
                      :fill="seg.fill"
                      :stroke="seg.stroke"
                      :stroke-width="seg.strokeWidth"
                      :stroke-dasharray="seg.dashArray"
                      :stroke-dashoffset="seg.dashOffset"
                      transform="rotate(-90 100 100)"
                    />
                    <text x="100" y="92" text-anchor="middle" class="pie-center-num">{{ platformStats.totalUsers.toLocaleString() }}</text>
                    <text x="100" y="112" text-anchor="middle" class="pie-center-label">服务用户</text>
                  </svg>
                  <div class="pie-legend">
                    <div v-for="(seg, i) in industryPieSegments" :key="i" class="pie-legend-item">
                      <span class="pie-legend-dot" :style="{ background: seg.color }"></span>
                      <span class="pie-legend-name">{{ seg.label }}</span>
                      <span class="pie-legend-value">{{ seg.value.toLocaleString() }}</span>
                      <span class="pie-legend-pct">{{ seg.pct }}%</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 求职转化漏斗 -->
            <div class="platform-card">
              <div class="platform-card-head">
                <div class="platform-card-title">
                  <RiseOutlined /> 求职转化漏斗
                </div>
                <span class="platform-card-sub">各阶段转化率（近 30 天）</span>
              </div>
              <div class="funnel-platform">
                <div v-for="(stage, i) in platformFunnel" :key="i" class="funnel-platform-row">
                  <div class="funnel-stage-label">{{ stage.label }}</div>
                  <div class="funnel-stage-bar-wrap">
                    <div class="funnel-stage-bar" :style="{ width: stage.pct + '%', background: stage.color }"></div>
                  </div>
                  <div class="funnel-stage-value">{{ stage.value.toLocaleString() }}</div>
                  <div class="funnel-stage-pct">{{ stage.pct }}%</div>
                  <div class="funnel-stage-conv" v-if="i > 0">
                    <ArrowRightOutlined />
                    {{ stage.convPct }}%
                  </div>
                </div>
              </div>
            </div>

            <div class="platform-grid-2">
              <!-- 热门岗位方向 -->
              <div class="platform-card">
                <div class="platform-card-head">
                  <div class="platform-card-title">
                    <StarOutlined /> 热门求职方向
                  </div>
                  <span class="platform-card-sub">投递量 Top 5</span>
                </div>
                <div class="hot-job-row" v-for="(job, i) in hotJobs" :key="i">
                  <span class="hot-job-rank" :class="'rank-' + (i + 1)">{{ i + 1 }}</span>
                  <span class="hot-job-name">{{ job.name }}</span>
                  <div class="hot-job-bar-wrap">
                    <div class="hot-job-bar" :style="{ width: job.pct + '%' }"></div>
                  </div>
                  <span class="hot-job-count">{{ job.count.toLocaleString() }}</span>
                </div>
              </div>

              <!-- 渠道效果对比 -->
              <div class="platform-card">
                <div class="platform-card-head">
                  <div class="platform-card-title">
                    <ThunderboltOutlined /> 渠道效果对比
                  </div>
                  <span class="platform-card-sub">各渠道 Offer 率</span>
                </div>
                <div class="channel-cmp-row" v-for="(ch, i) in channelComparison" :key="i">
                  <span class="channel-cmp-name">{{ ch.name }}</span>
                  <div class="channel-cmp-bar-wrap">
                    <div class="channel-cmp-bar" :style="{ width: ch.pct + '%', background: ch.color }"></div>
                  </div>
                  <span class="channel-cmp-pct">{{ ch.pct }}%</span>
                  <span class="channel-cmp-count">{{ ch.count }} 单</span>
                </div>
              </div>
            </div>

            <!-- 底部说明 -->
            <div class="platform-footer">
              <InfoCircleOutlined />
              <span>数据为平台模拟统计指标，仅用于展示产品功能效果</span>
              <span class="platform-update">最后更新：{{ platformUpdated }}</span>
            </div>
          </div>
        </Transition>
      </div>

      <!-- ==================== 详情面板 ==================== -->
      <Transition name="slide-fade">
        <aside class="detail-panel" v-if="selectedDelivery && viewMode === 'personal'" :key="selectedDelivery.id">
          <!-- 公司头部 -->
          <div class="dp-header">
            <div class="dp-header-left">
              <div class="dp-company-name">{{ selectedDelivery.company }}</div>
              <div class="dp-company-sub">{{ selectedDelivery.position }}</div>
            </div>
            <button class="dp-close" @click="selectedId = null" title="关闭 (Esc)">
              <CloseOutlined />
            </button>
          </div>

          <div class="dp-tags">
            <span class="kb-status" :class="getStatusClass(selectedDelivery.status)">
              <span class="status-dot" :style="{ background: getStatusColor(selectedDelivery.status) }"></span>
              {{ getStatusText(selectedDelivery.status) }}
            </span>
            <span class="kb-priority" :class="getPriorityClass(selectedDelivery.priority)">{{ getPriorityText(selectedDelivery.priority) }}优先级</span>
            <span class="dp-channel-tag">{{ getChannelLabel(selectedDelivery.channel) }}</span>
          </div>

          <div class="dp-meta-grid">
            <div class="dp-meta-item">
              <CalendarOutlined />
              <span>{{ formatDate(selectedDelivery.apply_date) }} 投递</span>
            </div>
            <div class="dp-meta-item" v-if="selectedDelivery.resume_version_id">
              <FileTextOutlined />
              <span>简历 v{{ selectedDelivery.resume_version_id }}</span>
            </div>
            <div class="dp-meta-item" v-if="parseHrContact(selectedDelivery.hr_contact)?.name">
              <UserOutlined />
              <span>{{ parseHrContact(selectedDelivery.hr_contact).name }}</span>
            </div>
          </div>

          <!-- Offer 详情（仅 Offer 状态显示） -->
          <div v-if="selectedDelivery.status === 'offer' && parseOfferDetail(selectedDelivery.offer_detail)" class="dp-offer-card">
            <div class="dp-offer-title">
              <TrophyOutlined /> Offer 详情
            </div>
            <div class="dp-offer-grid">
              <div v-if="parseOfferDetail(selectedDelivery.offer_detail).salary_base" class="dp-offer-item">
                <div class="dp-offer-label">基础薪资</div>
                <div class="dp-offer-value">{{ parseOfferDetail(selectedDelivery.offer_detail).salary_base }}</div>
              </div>
              <div v-if="parseOfferDetail(selectedDelivery.offer_detail).annual_bonus" class="dp-offer-item">
                <div class="dp-offer-label">年终奖</div>
                <div class="dp-offer-value">{{ parseOfferDetail(selectedDelivery.offer_detail).annual_bonus }}</div>
              </div>
              <div v-if="parseOfferDetail(selectedDelivery.offer_detail).stock" class="dp-offer-item">
                <div class="dp-offer-label">股票</div>
                <div class="dp-offer-value">{{ parseOfferDetail(selectedDelivery.offer_detail).stock }}</div>
              </div>
              <div v-if="parseOfferDetail(selectedDelivery.offer_detail).deadline" class="dp-offer-item">
                <div class="dp-offer-label">回复截止</div>
                <div class="dp-offer-value">{{ parseOfferDetail(selectedDelivery.offer_detail).deadline }}</div>
              </div>
            </div>
          </div>

          <!-- 面试时间线 -->
          <div class="dp-section-title">
            <span>面试时间线</span>
            <button class="dp-add-btn" @click="showRoundModal()"><PlusOutlined /> 新增</button>
          </div>
          <div class="tl-list">
            <!-- 投递节点 -->
            <div class="tl-item line-solid">
              <div class="tl-dot gray"></div>
              <div class="tl-content">
                <div class="tl-title">投递简历</div>
                <div class="tl-time">{{ formatDate(selectedDelivery.apply_date) }} · {{ getChannelLabel(selectedDelivery.channel) }}</div>
              </div>
            </div>
            <!-- 轮次节点 -->
            <div
              v-for="(round, idx) in sortedRounds"
              :key="round.id"
              class="tl-item"
              :class="getTimelineLineClass(round, idx)"
            >
              <div class="tl-dot" :class="getTimelineDotClass(round)"></div>
              <div class="tl-content">
                <div class="round-head">
                  <span class="tl-title">{{ getRoundTypeText(round.round_type) }}</span>
                  <a-dropdown :trigger="['click']">
                    <button class="kb-action kb-action-sm" @click.stop><EllipsisOutlined /></button>
                    <template #overlay>
                      <a-menu>
                        <a-menu-item @click="showRoundModal(round)">编辑</a-menu-item>
                        <a-menu-item danger @click="handleDeleteRound(round.id)">删除</a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
                </div>
                <div class="tl-time">
                  {{ round.interview_time ? formatDateTime(round.interview_time) : '时间未定' }}
                  <span v-if="round.result" :class="getRoundResultClass(round.result)"> · {{ getRoundResultText(round.result) }}</span>
                </div>
                <div class="tl-meta" v-if="round.interviewer_name">
                  <UserOutlined /> {{ round.interviewer_name }}{{ round.interviewer_title ? ' · ' + round.interviewer_title : '' }} · {{ getRoundFormatText(round.format) }}
                </div>
                <div class="tl-meta" v-if="round.question_summary">
                  <FileTextOutlined /> {{ round.question_summary }}
                </div>
                <span class="tl-tag" v-if="round.feedback">反馈：{{ round.feedback }}</span>
              </div>
            </div>
            <!-- 当前阶段提示 -->
            <div v-if="getNextInterview(selectedDelivery)" class="tl-item current line-dashed">
              <div class="tl-dot blue"></div>
              <div class="tl-content">
                <div class="tl-title">{{ getRoundTypeText(getNextInterview(selectedDelivery).round_type) }} · 待进行</div>
                <div class="tl-time">{{ formatDateTime(getNextInterview(selectedDelivery).interview_time) }}</div>
                <div class="tl-action-row">
                  <button class="dp-link-btn" @click="goToInterviewTraining(selectedDelivery)">
                    <ExportOutlined /> 去面试训练场练习
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- HR 反馈 -->
          <div class="dp-section-title">
            <span>HR 反馈</span>
            <button class="dp-add-btn" @click="showFeedbackModal()"><PlusOutlined /> 新增</button>
          </div>
          <div v-if="feedbacks.length > 0" class="fb-list">
            <div v-for="fb in feedbacks" :key="fb.id" class="fb-bubble">
              <div class="fb-text">{{ fb.summary }}</div>
              <div class="fb-meta">
                <span><MessageOutlined /> {{ formatDateTime(fb.contact_time) }} · {{ getMethodText(fb.method) }}</span>
                <button class="fb-del" @click="handleDeleteFeedback(fb.id)"><DeleteOutlined /></button>
              </div>
              <div class="fb-next" v-if="fb.next_action">
                <ArrowRightOutlined /> {{ fb.next_action }}
              </div>
            </div>
          </div>
          <div v-else class="empty-hint">暂无 HR 反馈记录</div>

          <!-- AI 策略洞察 -->
          <div class="strategy-card">
            <div class="strategy-title">
              <RiseOutlined /> 求职策略洞察
            </div>

              <!-- 投递漏斗 SVG -->
              <div class="strategy-sub">投递漏斗</div>
              <div class="funnel-svg-wrap">
                <svg class="funnel-svg" viewBox="0 0 280 100" preserveAspectRatio="none">
                  <polygon
                    v-for="(seg, i) in funnelSegments"
                    :key="i"
                    :points="seg.points"
                    :fill="seg.color"
                    :opacity="0.85"
                  />
                </svg>
                <div class="funnel-labels">
                  <div v-for="(seg, i) in funnelSegments" :key="i" class="funnel-label-item">
                    <span class="funnel-label-name">{{ seg.label }}</span>
                    <span class="funnel-label-val">{{ seg.value }} <span class="funnel-label-pct">({{ seg.pct }}%)</span></span>
                  </div>
                </div>
              </div>

              <!-- 渠道分布 -->
              <div class="strategy-sub">渠道分布</div>
              <div class="bar-row" v-for="ch in channelStats" :key="ch.label">
                <span class="bar-label">{{ ch.label }}<span class="mono">{{ ch.value }}</span></span>
                <div class="bar-track"><div class="bar-fill" :class="ch.cls" :style="{ width: ch.pct + '%' }"></div></div>
              </div>

              <!-- 平均周期 -->
              <div class="strategy-sub">平均周期</div>
              <div class="cycle-stat">
                <span class="cycle-num">{{ avgCycle }}</span>
                <span class="cycle-unit">天 · 从投递到 Offer</span>
              </div>
          </div>
        </aside>
      </Transition>

      <!-- 未选中时的占位（仅 personal 视图） -->
      <aside class="detail-panel detail-empty" v-if="!selectedDelivery && viewMode === 'personal'">
        <div class="detail-empty-icon">
          <AimOutlined />
        </div>
        <div class="detail-empty-title">点击投递查看详情</div>
        <div class="detail-empty-desc">选择左侧任意投递记录，查看面试时间线、HR 反馈与策略洞察</div>
        <div class="detail-empty-shortcuts">
          <div class="shortcut-row"><kbd>J</kbd><kbd>K</kbd> 上下导航</div>
          <div class="shortcut-row"><kbd>N</kbd> 新增投递</div>
          <div class="shortcut-row"><kbd>Esc</kbd> 关闭详情</div>
        </div>
      </aside>
    </div>

    <!-- ==================== 全局空状态（无任何投递） ==================== -->
    <div v-else class="empty-state empty-state-full">
      <div class="empty-illustration">
        <SendOutlined />
      </div>
      <div class="empty-title">开始追踪你的投递</div>
      <div class="empty-desc">使用顶部 Quick Capture 快速创建投递，或点击「新增投递」填写完整信息</div>
      <button class="btn-primary-capsule" @click="showCreateModal">
        <PlusOutlined /> 新增第一条投递
      </button>
    </div>

    <!-- ==================== 新增投递弹窗 ==================== -->
    <a-modal
      v-model:open="createModalVisible"
      title="新增投递"
      width="640px"
      @ok="handleCreate"
      @cancel="resetCreateForm"
    >
      <a-form :model="createForm" :rules="createRules" layout="vertical" ref="createFormRef">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="公司" name="company">
              <a-input v-model:value="createForm.company" placeholder="请输入公司名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="职位" name="position">
              <a-input v-model:value="createForm.position" placeholder="请输入职位名称" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="8">
            <a-form-item label="投递渠道" name="channel">
              <a-select v-model:value="createForm.channel" placeholder="请选择">
                <a-select-option v-for="c in channelOptions" :key="c.value" :value="c.value">{{ c.label }}</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="投递日期" name="apply_date">
              <a-date-picker v-model:value="createForm.apply_date" placeholder="选择日期" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="优先级" name="priority">
              <a-select v-model:value="createForm.priority" placeholder="请选择">
                <a-select-option value="high">高</a-select-option>
                <a-select-option value="medium">中</a-select-option>
                <a-select-option value="low">低</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="JD 链接" name="jd_text">
          <a-input v-model:value="createForm.jd_text" placeholder="职位 JD 链接（可选）" />
        </a-form-item>
        <a-form-item label="JD 内容" name="jd_content">
          <a-textarea v-model:value="createForm.remark" placeholder="职位 JD 描述（可选，用于联动面试训练场）" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- ==================== 轮次表单弹窗 ==================== -->
    <a-modal
      v-model:open="roundModalVisible"
      :title="editingRound ? '编辑轮次' : '新增轮次'"
      width="560px"
      @ok="handleSaveRound"
    >
      <a-form :model="roundForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="轮次类型" name="round_type">
              <a-select v-model:value="roundForm.round_type" placeholder="请选择">
                <a-select-option value="written_test">笔试</a-select-option>
                <a-select-option value="first_tech">技术一面</a-select-option>
                <a-select-option value="second_tech">技术二面</a-select-option>
                <a-select-option value="third_tech">技术三面</a-select-option>
                <a-select-option value="cross">交叉面</a-select-option>
                <a-select-option value="hr">HR 面</a-select-option>
                <a-select-option value="additional">加面</a-select-option>
                <a-select-option value="final">终面</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="面试形式" name="format">
              <a-select v-model:value="roundForm.format" placeholder="请选择">
                <a-select-option value="onsite">现场</a-select-option>
                <a-select-option value="video">视频</a-select-option>
                <a-select-option value="phone">电话</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="面试时间" name="interview_time">
              <a-date-picker v-model:value="roundForm.interview_time" show-time placeholder="选择时间" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="结果" name="result">
              <a-select v-model:value="roundForm.result" placeholder="请选择">
                <a-select-option value="pass">通过</a-select-option>
                <a-select-option value="pending">待定</a-select-option>
                <a-select-option value="rejected">未通过</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="面试官姓名" name="interviewer_name">
              <a-input v-model:value="roundForm.interviewer_name" placeholder="可选" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="面试官头衔" name="interviewer_title">
              <a-input v-model:value="roundForm.interviewer_title" placeholder="可选" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="问题摘要" name="question_summary">
          <a-textarea v-model:value="roundForm.question_summary" :rows="2" placeholder="可选" />
        </a-form-item>
        <a-form-item label="反馈" name="feedback">
          <a-textarea v-model:value="roundForm.feedback" :rows="2" placeholder="可选" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- ==================== 反馈表单弹窗 ==================== -->
    <a-modal
      v-model:open="feedbackModalVisible"
      title="新增 HR 反馈"
      width="480px"
      @ok="handleSaveFeedback"
    >
      <a-form :model="feedbackForm" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="联系时间" name="contact_time" required>
              <a-date-picker v-model:value="feedbackForm.contact_time" show-time placeholder="选择时间" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="联系方式" name="method" required>
              <a-select v-model:value="feedbackForm.method" placeholder="请选择">
                <a-select-option value="wechat">微信</a-select-option>
                <a-select-option value="phone">电话</a-select-option>
                <a-select-option value="email">邮件</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="反馈摘要" name="summary" required>
          <a-textarea v-model:value="feedbackForm.summary" :rows="3" placeholder="HR 反馈内容" />
        </a-form-item>
        <a-form-item label="下一步行动" name="next_action">
          <a-textarea v-model:value="feedbackForm.next_action" :rows="2" placeholder="可选" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useDeliveryStore } from '@/stores/delivery'
import { message, Modal } from 'ant-design-vue'
import type { FormInstance } from 'ant-design-vue'
import type {
  Delivery,
  DeliveryStatus,
  DeliveryRound,
  CreateDeliveryRequest,
  CreateRoundRequest,
  CreateFeedbackRequest,
} from '@/types/models'
import dayjs, { Dayjs } from 'dayjs'
// 使用 @ant-design/icons-vue，不再使用 lucide-vue-next
import {
  SendOutlined,
  TeamOutlined,
  TrophyOutlined,
  CloseCircleOutlined,
  TableOutlined,
  PieChartOutlined,
  SearchOutlined,
  PlusOutlined,
  RightOutlined,
  RightCircleOutlined,
  EllipsisOutlined,
  CalendarOutlined,
  FileTextOutlined,
  UserOutlined,
  MessageOutlined,
  ArrowRightOutlined,
  CloseOutlined,
  ExportOutlined,
  RiseOutlined,
  ThunderboltOutlined,
  BankOutlined,
  SolutionOutlined,
  DeleteOutlined,
  InboxOutlined,
  AimOutlined,
  StarOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons-vue'

const router = useRouter()
const deliveryStore = useDeliveryStore()

// ==================== 视图与筛选 ====================
const viewMode = ref<'personal' | 'platform'>('personal')
const searchKeyword = ref('')
const hoveredId = ref<number | null>(null)

const filters = reactive({
  status: '' as string,
  channel: '' as string,
  priority: '' as string,
})

// 状态枚举（6 列，对齐后端 + 原型配色）
const statusList = [
  { label: '待响应', value: 'pending' as DeliveryStatus, dotColor: '#aeaeb2' },
  { label: '笔试中', value: 'written_test' as DeliveryStatus, dotColor: '#ff9500' },
  { label: '面试中', value: 'interview' as DeliveryStatus, dotColor: '#007aff' },
  { label: '待Offer', value: 'waiting_offer' as DeliveryStatus, dotColor: '#af52de' },
  { label: '已Offer', value: 'offer' as DeliveryStatus, dotColor: '#34c759' },
  { label: '已拒绝', value: 'rejected' as DeliveryStatus, dotColor: '#ff3b30' },
]

const channelOptions = [
  { value: 'boss', label: 'BOSS直聘' },
  { value: 'official', label: '官网' },
  { value: 'referral', label: '内推' },
  { value: 'campus', label: '校园招聘' },
  { value: 'headhunt', label: '猎头' },
  { value: 'other', label: '其他' },
]

// 状态合法流转
const transitionMap: Record<DeliveryStatus, DeliveryStatus[]> = {
  pending: ['written_test', 'interview', 'rejected'],
  written_test: ['interview', 'waiting_offer', 'rejected'],
  interview: ['waiting_offer', 'offer', 'rejected'],
  waiting_offer: ['offer', 'rejected'],
  offer: [],
  rejected: ['interview'],
}

const getAvailableTransitions = (status: DeliveryStatus) => {
  return transitionMap[status]?.map((v) => ({ value: v, label: getStatusText(v) })) || []
}

// 本地筛选 + 搜索
const filteredDeliveries = computed(() => {
  return deliveryStore.deliveries.filter((d) => {
    if (filters.status && d.status !== filters.status) return false
    if (filters.channel && d.channel !== filters.channel) return false
    if (filters.priority && d.priority !== filters.priority) return false
    if (searchKeyword.value) {
      const kw = searchKeyword.value.toLowerCase()
      if (!d.company.toLowerCase().includes(kw) && !d.position.toLowerCase().includes(kw)) return false
    }
    return true
  })
})

const hasActiveFilter = computed(() =>
  !!(filters.status || filters.channel || filters.priority || searchKeyword.value)
)

const clearFilters = () => {
  filters.status = ''
  filters.channel = ''
  filters.priority = ''
  searchKeyword.value = ''
}

// ==================== 统计数据 + Sparkline ====================
const stats = computed(() => {
  const total = deliveryStore.deliveries.length
  const inProgress = deliveryStore.deliveries.filter(
    (d) => d.status === 'interview' || d.status === 'written_test'
  ).length
  const offerCount = deliveryStore.deliveries.filter((d) => d.status === 'offer').length
  const rejected = deliveryStore.deliveries.filter((d) => d.status === 'rejected').length
  return { total, inProgress, offerCount, rejected }
})

const rejectRate = computed(() => {
  if (stats.value.total === 0) return 0
  return Math.round((stats.value.rejected / stats.value.total) * 100)
})

const weekCount = computed(() => {
  const weekAgo = dayjs().subtract(7, 'day')
  return deliveryStore.deliveries.filter((d) => dayjs(d.apply_date).isAfter(weekAgo)).length
})

const offerWeekCount = computed(() => {
  const weekAgo = dayjs().subtract(7, 'day')
  return deliveryStore.deliveries.filter(
    (d) => d.status === 'offer' && dayjs(d.updated_at).isAfter(weekAgo)
  ).length
})

// 生成 7 天 sparkline 数据（基于 apply_date）
const generateSparkline = (predicate: (d: Delivery) => boolean): string => {
  const points: number[] = []
  for (let i = 6; i >= 0; i--) {
    const day = dayjs().subtract(i, 'day')
    const next = day.add(1, 'day')
    const count = deliveryStore.deliveries.filter(
      (d) => predicate(d) && dayjs(d.apply_date).isAfter(day.subtract(1, 'ms')) && dayjs(d.apply_date).isBefore(next)
    ).length
    points.push(count)
  }
  const max = Math.max(...points, 1)
  return points.map((p, i) => `${(i / 6) * 80},${24 - (p / max) * 22 - 1}`).join(' ')
}

// 统计卡片：1:1 对齐设计稿配色（蓝/紫/绿/红）
const statCards = computed(() => [
  {
    key: 'total',
    label: '投递总数',
    value: stats.value.total,
    change: `↑ +${weekCount.value} 本周`,
    changeClass: 'up',
    icon: SendOutlined,
    iconBg: 'var(--brand-50)',
    iconColor: 'var(--primary)',
    accentClass: 'accent-blue',
    sparkline: generateSparkline(() => true),
    sparkColor: '#007aff',
  },
  {
    key: 'progress',
    label: '面试中',
    value: stats.value.inProgress,
    change: `↑ 进行中`,
    changeClass: 'up',
    icon: TeamOutlined,
    iconBg: 'rgba(175,82,222,0.14)',
    iconColor: '#8e3db5',
    accentClass: 'accent-purple',
    sparkline: generateSparkline((d) => d.status === 'interview' || d.status === 'written_test'),
    sparkColor: '#af52de',
  },
  {
    key: 'offer',
    label: 'Offer',
    value: stats.value.offerCount,
    change: `↑ +${offerWeekCount.value} 本月`,
    changeClass: 'up',
    icon: TrophyOutlined,
    iconBg: 'var(--state-success-surface)',
    iconColor: 'var(--success)',
    accentClass: 'accent-green',
    sparkline: generateSparkline((d) => d.status === 'offer'),
    sparkColor: '#34c759',
  },
  {
    key: 'rejected',
    label: '拒绝',
    value: stats.value.rejected,
    change: `↓ 转化率 ${rejectRate.value}%`,
    changeClass: 'down',
    icon: CloseCircleOutlined,
    iconBg: 'var(--state-error-surface)',
    iconColor: 'var(--state-error)',
    accentClass: 'accent-red',
    sparkline: generateSparkline((d) => d.status === 'rejected'),
    sparkColor: '#ff3b30',
  },
])

// ==================== 选中与详情 ====================
const selectedId = ref<number | null>(null)

const selectedDelivery = computed(() => {
  if (!selectedId.value) return null
  return deliveryStore.deliveries.find((d) => d.id === selectedId.value) || deliveryStore.currentDelivery
})

const rounds = computed(() => deliveryStore.rounds)
const feedbacks = computed(() => deliveryStore.feedbacks)

const sortedRounds = computed(() => {
  return [...deliveryStore.rounds]
    .filter((r) => r.result !== 'pending' || r.interview_time)
    .sort((a, b) => {
      const ta = a.interview_time ? new Date(a.interview_time).getTime() : 0
      const tb = b.interview_time ? new Date(b.interview_time).getTime() : 0
      return ta - tb
    })
})

const selectDelivery = async (delivery: Delivery) => {
  selectedId.value = delivery.id
  await deliveryStore.fetchDelivery(delivery.id)
}

// ==================== 进度点（5 个：笔试/一面/二面/HR面/终面）====================
const progressSlots = ['written_test', 'first_tech', 'second_tech', 'hr', 'final']

const getProgressDots = (delivery: Delivery): string[] => {
  const dots: string[] = []
  for (const slot of progressSlots) {
    const round = deliveryStore.rounds.find((r) => r.round_type === slot)
    if (!round) {
      dots.push('')
    } else if (round.result === 'pass') {
      dots.push('done')
    } else if (round.result === 'rejected') {
      dots.push('fail')
    } else {
      dots.push('current')
    }
  }
  // 无轮次数据时，根据状态推断
  if (deliveryStore.rounds.length === 0) {
    const statusProgress: Record<string, number> = {
      pending: 0,
      written_test: 1,
      interview: 2,
      waiting_offer: 4,
      offer: 5,
      rejected: 0,
    }
    const currentStep = statusProgress[delivery.status] ?? 0
    return progressSlots.map((_, i) => {
      if (i < currentStep) return 'done'
      if (i === currentStep && delivery.status !== 'offer' && delivery.status !== 'pending') return 'current'
      if (delivery.status === 'offer') return 'done'
      return ''
    })
  }
  return dots
}

const getProgressTooltip = (delivery: Delivery): string => {
  const labels = ['笔试', '一面', '二面', 'HR面', '终面']
  const dots = getProgressDots(delivery)
  return labels.map((l, i) => {
    const d = dots[i]
    if (d === 'done') return `${l}✓`
    if (d === 'current') return `${l}○`
    if (d === 'fail') return `${l}✗`
    return `${l}—`
  }).join(' ')
}

// ==================== 下次面试 / HR 最新反馈 ====================
const getNextInterview = (delivery: Delivery): DeliveryRound | null => {
  if (delivery.id !== selectedId.value && deliveryStore.rounds.length === 0) return null
  if (delivery.id !== selectedId.value) return null
  const now = Date.now()
  return sortedRounds.value.find(
    (r) => r.result === 'pending' && r.interview_time && new Date(r.interview_time).getTime() > now
  ) || null
}

const getLatestFeedback = (deliveryId: number) => {
  if (deliveryId !== selectedId.value) return null
  return feedbacks.value[0] || null
}

// ==================== 漏斗与渠道统计 ====================
const funnelData = computed(() => {
  const total = stats.value.total || 1
  const written = deliveryStore.deliveries.filter(
    (d) => ['written_test', 'interview', 'waiting_offer', 'offer'].includes(d.status)
  ).length
  const interview = deliveryStore.deliveries.filter(
    (d) => ['interview', 'waiting_offer', 'offer'].includes(d.status)
  ).length
  const waiting = deliveryStore.deliveries.filter(
    (d) => ['waiting_offer', 'offer'].includes(d.status)
  ).length
  const offer = stats.value.offerCount
  return [
    { label: '投递', value: stats.value.total, pct: 100 },
    { label: '笔试', value: written, pct: Math.round((written / total) * 100) },
    { label: '面试', value: interview, pct: Math.round((interview / total) * 100) },
    { label: '待Offer', value: waiting, pct: Math.round((waiting / total) * 100) },
    { label: 'Offer', value: offer, pct: Math.round((offer / total) * 100) },
  ]
})

// SVG 漏斗图分段
const funnelSegments = computed(() => {
  const data = funnelData.value
  const maxW = 280
  const segH = 18
  const gap = 2
  const colors = ['#007aff', '#5e5ce6', '#af52de', '#ff9500', '#34c759']
  const maxVal = data[0].value || 1
  return data.map((seg, i) => {
    const w = Math.max(20, (seg.value / maxVal) * maxW)
    const x1 = (maxW - w) / 2
    const x2 = x1 + w
    const y1 = i * (segH + gap)
    const y2 = y1 + segH
    const prevW = i === 0 ? maxW : Math.max(20, (data[i - 1].value / maxVal) * maxW)
    const prevX1 = (maxW - prevW) / 2
    const prevX2 = prevX1 + prevW
    return {
      label: seg.label,
      value: seg.value,
      pct: seg.pct,
      color: colors[i],
      points: `${prevX1},${y1} ${prevX2},${y1} ${x2},${y2} ${x1},${y2}`,
    }
  })
})

const channelStats = computed(() => {
  const counts: Record<string, number> = {}
  deliveryStore.deliveries.forEach((d) => {
    counts[d.channel] = (counts[d.channel] || 0) + 1
  })
  const max = Math.max(...Object.values(counts), 1)
  const clsMap: Record<string, string> = { boss: 'boss', official: 'web', referral: 'ref', campus: 'you', headhunt: 'avg', other: 'avg' }
  return Object.entries(counts).map(([ch, val]) => ({
    label: getChannelLabel(ch),
    value: val,
    pct: Math.round((val / max) * 100),
    cls: clsMap[ch] || 'avg',
  }))
})

const avgCycle = computed(() => {
  const offered = deliveryStore.deliveries.filter((d) => d.status === 'offer')
  if (offered.length === 0) return 0
  const days = offered.map((d) => {
    const start = dayjs(d.apply_date)
    const end = dayjs(d.updated_at)
    return end.diff(start, 'day')
  })
  return Math.round(days.reduce((a, b) => a + b, 0) / days.length)
})

// ==================== Quick Capture（顶部快捷输入）====================
const quickCapture = ref('')
const quickCaptureParsed = ref<{ company?: string; position?: string }>({})

const handlePaste = (e: ClipboardEvent) => {
  const text = e.clipboardData?.getData('text') || ''
  // 简单 JD 解析：尝试匹配「公司：xxx」「岗位：xxx」
  const companyMatch = text.match(/公司[：:]\s*(.+)/)
  const positionMatch = text.match(/岗位[：:]\s*(.+)/) || text.match(/职位[：:]\s*(.+)/)
  quickCaptureParsed.value = {
    company: companyMatch?.[1]?.trim(),
    position: positionMatch?.[1]?.trim(),
  }
}

const handleQuickCapture = () => {
  const text = quickCapture.value.trim()
  if (!text) return
  const parsed = quickCaptureParsed.value
  if (parsed.company && parsed.position) {
    createFromQuickCapture(parsed.company, parsed.position, text)
    return
  }
  if (text.includes('/')) {
    const [company, position] = text.split('/').map((s) => s.trim())
    if (company && position) {
      createFromQuickCapture(company, position, '')
      return
    }
  }
  createForm.company = text
  createForm.jd_text = text.length > 20 ? text : ''
  createModalVisible.value = true
  quickCapture.value = ''
  quickCaptureParsed.value = {}
}

const createFromQuickCapture = async (company: string, position: string, jdText: string) => {
  const data: CreateDeliveryRequest = {
    company,
    position,
    channel: 'other',
    apply_date: dayjs().format('YYYY-MM-DD'),
    priority: 'medium',
    jd_text: jdText,
  }
  await deliveryStore.createDelivery(data)
  quickCapture.value = ''
  quickCaptureParsed.value = {}
}

// ==================== 键盘快捷键 ====================
const handleKeydown = (e: KeyboardEvent) => {
  const target = e.target as HTMLElement
  if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA') return
  if (target.isContentEditable) return

  const list = filteredDeliveries.value
  const currentIdx = list.findIndex((d) => d.id === selectedId.value)

  if (e.key === 'j' || e.key === 'J') {
    e.preventDefault()
    const next = list[Math.min(currentIdx + 1, list.length - 1)] || list[0]
    if (next) selectDelivery(next)
  } else if (e.key === 'k' || e.key === 'K') {
    e.preventDefault()
    const prev = list[Math.max(currentIdx - 1, 0)]
    if (prev) selectDelivery(prev)
  } else if (e.key === 'n' || e.key === 'N') {
    e.preventDefault()
    showCreateModal()
  } else if (e.key === 'Escape') {
    selectedId.value = null
  }
}

// ==================== 新增投递 ====================
const createModalVisible = ref(false)
const createFormRef = ref<FormInstance>()
const createForm = reactive<CreateDeliveryRequest & { apply_date?: Dayjs }>({
  company: '',
  position: '',
  channel: 'boss',
  apply_date: undefined,
  priority: 'medium',
  jd_text: '',
  remark: '',
  resume_version_id: undefined,
})

const createRules = {
  company: [{ required: true, message: '请输入公司名称', trigger: 'blur' }],
  position: [{ required: true, message: '请输入职位名称', trigger: 'blur' }],
  channel: [{ required: true, message: '请选择投递渠道', trigger: 'change' }],
  apply_date: [{ required: true, message: '请选择投递日期', trigger: 'change' }],
}

const showCreateModal = () => {
  createModalVisible.value = true
}

const handleCreate = async () => {
  try {
    await createFormRef.value?.validateFields()
    const data: CreateDeliveryRequest = {
      company: createForm.company,
      position: createForm.position,
      channel: createForm.channel as CreateDeliveryRequest['channel'],
      apply_date: createForm.apply_date?.format('YYYY-MM-DD') || '',
      priority: createForm.priority as CreateDeliveryRequest['priority'],
      jd_text: createForm.jd_text,
      remark: createForm.remark,
      resume_version_id: createForm.resume_version_id,
    }
    await deliveryStore.createDelivery(data)
    createModalVisible.value = false
    resetCreateForm()
  } catch (error) {
    console.error('验证失败:', error)
  }
}

const resetCreateForm = () => {
  createFormRef.value?.resetFields()
  Object.assign(createForm, {
    company: '',
    position: '',
    channel: 'boss',
    apply_date: undefined,
    priority: 'medium',
    jd_text: '',
    remark: '',
    resume_version_id: undefined,
  })
}

// ==================== 状态变更与删除 ====================
const handleStatusChange = async (delivery: Delivery, status: DeliveryStatus) => {
  await deliveryStore.changeStatus(delivery.id, { status })
}

const handleDelete = (id: number) => {
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除这条投递记录吗？',
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      await deliveryStore.deleteDelivery(id)
      if (selectedId.value === id) selectedId.value = null
    },
  })
}

// ==================== 轮次操作 ====================
const roundModalVisible = ref(false)
const editingRound = ref<DeliveryRound | null>(null)
const roundForm = reactive<{
  round_type: string
  interview_time: Dayjs | null
  format: string
  interviewer_name: string
  interviewer_title: string
  question_summary: string
  feedback: string
  result: string
}>({
  round_type: 'first_tech',
  interview_time: null,
  format: 'video',
  interviewer_name: '',
  interviewer_title: '',
  question_summary: '',
  feedback: '',
  result: 'pending',
})

const showRoundModal = (round?: DeliveryRound) => {
  if (round) {
    editingRound.value = round
    Object.assign(roundForm, {
      round_type: round.round_type,
      interview_time: round.interview_time ? dayjs(round.interview_time) : null,
      format: round.format,
      interviewer_name: round.interviewer_name,
      interviewer_title: round.interviewer_title,
      question_summary: round.question_summary,
      feedback: round.feedback,
      result: round.result,
    })
  } else {
    editingRound.value = null
    Object.assign(roundForm, {
      round_type: 'first_tech',
      interview_time: null,
      format: 'video',
      interviewer_name: '',
      interviewer_title: '',
      question_summary: '',
      feedback: '',
      result: 'pending',
    })
  }
  roundModalVisible.value = true
}

const showRoundModalFor = (delivery: Delivery) => {
  selectDelivery(delivery).then(() => showRoundModal())
}

const handleSaveRound = async () => {
  if (!selectedDelivery.value) return
  const deliveryId = selectedDelivery.value.id
  const data: CreateRoundRequest = {
    round_type: roundForm.round_type,
    interview_time: roundForm.interview_time?.format('YYYY-MM-DD HH:mm:ss') || '',
    format: roundForm.format,
    interviewer_name: roundForm.interviewer_name,
    interviewer_title: roundForm.interviewer_title,
    question_summary: roundForm.question_summary,
    feedback: roundForm.feedback,
    result: roundForm.result,
  }
  if (editingRound.value) {
    await deliveryStore.updateRound(deliveryId, editingRound.value.id, data)
  } else {
    await deliveryStore.createRound(deliveryId, data)
  }
  roundModalVisible.value = false
}

const handleDeleteRound = (roundId: number) => {
  if (!selectedDelivery.value) return
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该面试轮次吗？',
    onOk: async () => {
      await deliveryStore.deleteRound(selectedDelivery.value!.id, roundId)
    },
  })
}

// ==================== 反馈操作 ====================
const feedbackModalVisible = ref(false)
const feedbackForm = reactive<{
  contact_time: Dayjs | null
  method: string
  summary: string
  next_action: string
}>({
  contact_time: null,
  method: 'wechat',
  summary: '',
  next_action: '',
})

const showFeedbackModal = () => {
  Object.assign(feedbackForm, {
    contact_time: dayjs(),
    method: 'wechat',
    summary: '',
    next_action: '',
  })
  feedbackModalVisible.value = true
}

const showFeedbackModalFor = (delivery: Delivery) => {
  selectDelivery(delivery).then(() => showFeedbackModal())
}

const handleSaveFeedback = async () => {
  if (!selectedDelivery.value) return
  if (!feedbackForm.contact_time || !feedbackForm.summary) {
    message.warning('请填写联系时间和反馈摘要')
    return
  }
  const data: CreateFeedbackRequest = {
    contact_time: feedbackForm.contact_time.format('YYYY-MM-DD HH:mm'),
    method: feedbackForm.method,
    summary: feedbackForm.summary,
    next_action: feedbackForm.next_action,
  }
  await deliveryStore.createFeedback(selectedDelivery.value.id, data)
  feedbackModalVisible.value = false
}

const handleDeleteFeedback = (feedbackId: number) => {
  if (!selectedDelivery.value) return
  Modal.confirm({
    title: '确认删除',
    content: '确定要删除该反馈吗？',
    onOk: async () => {
      await deliveryStore.deleteFeedback(selectedDelivery.value!.id, feedbackId)
    },
  })
}

// ==================== 跳转面试训练场 ====================
const goToInterviewTraining = (delivery: Delivery) => {
  router.push({
    path: '/interview',
    query: {
      company: delivery.company,
      position: delivery.position,
      jd: delivery.jd_text || '',
    },
  })
}

// ==================== 辅助函数 ====================
const handleFilter = async () => {
  await deliveryStore.fetchDeliveries({
    status: filters.status || undefined,
    channel: filters.channel || undefined,
  })
}

const getStatusText = (status: string): string => {
  return statusList.find((s) => s.value === status)?.label || status
}

const getStatusClass = (status: string): string => {
  const map: Record<string, string> = {
    pending: 'gray',
    written_test: 'orange',
    interview: 'blue',
    waiting_offer: 'purple',
    offer: 'green',
    rejected: 'red',
  }
  return map[status] || 'gray'
}

const getStatusColor = (status: string): string => {
  const map: Record<string, string> = {
    pending: '#aeaeb2',
    written_test: '#ff9500',
    interview: '#007aff',
    waiting_offer: '#af52de',
    offer: '#34c759',
    rejected: '#ff3b30',
  }
  return map[status] || '#aeaeb2'
}

const getChannelLabel = (channel: string): string => {
  return channelOptions.find((c) => c.value === channel)?.label || channel || '-'
}

const getChannelClass = (channel: string): string => {
  const map: Record<string, string> = {
    boss: 'ch-boss',
    official: 'ch-official',
    referral: 'ch-referral',
    campus: 'ch-campus',
    headhunt: 'ch-headhunt',
    other: 'ch-other',
  }
  return map[channel] || 'ch-other'
}

const getPriorityText = (priority: string): string => {
  const map: Record<string, string> = { high: '高', medium: '中', low: '低' }
  return map[priority] || priority || '-'
}

const getPriorityClass = (priority: string): string => {
  const map: Record<string, string> = { high: 'high', medium: 'mid', low: 'low' }
  return map[priority] || 'low'
}

const getRoundResultText = (result?: string): string => {
  if (!result) return '待定'
  const map: Record<string, string> = { pass: '通过', pending: '待定', rejected: '未通过' }
  return map[result] || result
}

const getRoundResultClass = (result?: string): string => {
  if (result === 'pass') return 'result-pass'
  if (result === 'rejected') return 'result-fail'
  return 'result-pending'
}

const getRoundTypeText = (type?: string): string => {
  const map: Record<string, string> = {
    written_test: '笔试',
    first_tech: '技术一面',
    second_tech: '技术二面',
    third_tech: '技术三面',
    cross: '交叉面',
    hr: 'HR 面',
    additional: '加面',
    final: '终面',
  }
  return map[type] || type || '未知'
}

const getRoundFormatText = (format?: string): string => {
  const map: Record<string, string> = { onsite: '现场', video: '视频', phone: '电话' }
  return map[format] || format || '-'
}

const getMethodText = (method?: string): string => {
  const map: Record<string, string> = { wechat: '微信', phone: '电话', email: '邮件' }
  return map[method] || method || '-'
}

const getTimelineLineClass = (round: DeliveryRound, idx: number): string => {
  if (idx === sortedRounds.value.length - 1) return 'line-none'
  return round.result === 'pass' ? 'line-solid' : 'line-dashed'
}

const getTimelineDotClass = (round: DeliveryRound): string => {
  if (round.result === 'pass') return 'green'
  if (round.result === 'rejected') return 'red'
  return 'blue'
}

const parseHrContact = (json: string): { name?: string; wechat?: string; phone?: string; email?: string } => {
  if (!json) return {}
  try {
    return JSON.parse(json)
  } catch {
    return {}
  }
}

const parseOfferDetail = (json: string): {
  salary_base?: string
  annual_bonus?: string
  stock?: string
  benefits?: string
  deadline?: string
} => {
  if (!json) return {}
  try {
    return JSON.parse(json)
  } catch {
    return {}
  }
}

const truncate = (text: string, len: number): string => {
  if (!text) return ''
  return text.length > len ? text.slice(0, len) + '…' : text
}

const formatDate = (date: string): string => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

const formatShortDateTime = (date: string): string => {
  if (!date) return '-'
  return dayjs(date).format('MM-DD HH:mm')
}

const formatDateTime = (date: string): string => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// ==================== 平台数据视图（按日期种子稳定） ====================
const seedRand = (seed: number) => {
  let s = seed
  return () => {
    s = (s * 9301 + 49297) % 233280
    return s / 233280
  }
}

const todaySeed = computed(() => {
  const d = new Date()
  return d.getFullYear() * 10000 + (d.getMonth() + 1) * 100 + d.getDate()
})

// 平台核心指标（基于种子稳定生成，避免每次刷新跳变）
const platformStats = computed(() => {
  const rand = seedRand(todaySeed.value)
  const totalUsers = 12800 + Math.floor(rand() * 3200)
  const totalDeliveries = 86000 + Math.floor(rand() * 8000)
  const successRate = 38 + Math.floor(rand() * 12)
  const avgCycleDays = 28 + Math.floor(rand() * 12)
  return { totalUsers, totalDeliveries, successRate, avgCycleDays }
})

// 投递结果分布饼图分段
const pieSegments = computed(() => {
  const rand = seedRand(todaySeed.value + 1)
  const offer = 4200 + Math.floor(rand() * 600)
  const interview = 8600 + Math.floor(rand() * 1200)
  const rejected = 12000 + Math.floor(rand() * 2000)
  const pending = platformStats.value.totalDeliveries - offer - interview - rejected
  const total = platformStats.value.totalDeliveries
  const segs = [
    { label: '已 Offer', value: offer, color: '#34c759' },
    { label: '面试中', value: interview, color: '#007aff' },
    { label: '已拒绝', value: rejected, color: '#ff3b30' },
    { label: '待响应', value: pending, color: '#aeaeb2' },
  ]
  const circumference = 2 * Math.PI * 70
  let accOffset = 0
  return segs.map((seg) => {
    const pct = total > 0 ? Math.round((seg.value / total) * 100) : 0
    const dashLen = (seg.value / total) * circumference
    const segData = {
      ...seg,
      pct,
      fill: 'none',
      stroke: seg.color,
      strokeWidth: 28,
      dashArray: `${dashLen} ${circumference - dashLen}`,
      dashOffset: -accOffset,
    }
    accOffset += dashLen
    return segData
  })
})

// 服务行业分布饼图分段
const industryPieSegments = computed(() => {
  const rand = seedRand(todaySeed.value + 2)
  const total = platformStats.value.totalUsers
  const data = [
    { label: '互联网', value: 4200 + Math.floor(rand() * 800), color: '#007aff' },
    { label: '金融', value: 2800 + Math.floor(rand() * 600), color: '#34c759' },
    { label: '制造', value: 2200 + Math.floor(rand() * 500), color: '#ff9500' },
    { label: '教育', value: 1800 + Math.floor(rand() * 400), color: '#5856d6' },
    { label: '医疗', value: 1400 + Math.floor(rand() * 300), color: '#af52de' },
    { label: '其他', value: total - 12400, color: '#aeaeb2' },
  ]
  const circumference = 2 * Math.PI * 70
  let accOffset = 0
  return data.map((seg) => {
    const pct = total > 0 ? Math.round((seg.value / total) * 100) : 0
    const dashLen = (seg.value / total) * circumference
    const segData = {
      ...seg,
      pct,
      fill: 'none',
      stroke: seg.color,
      strokeWidth: 28,
      dashArray: `${dashLen} ${circumference - dashLen}`,
      dashOffset: -accOffset,
    }
    accOffset += dashLen
    return segData
  })
})

// 求职转化漏斗（平台）
const platformFunnel = computed(() => {
  const rand = seedRand(todaySeed.value + 3)
  const total = platformStats.value.totalDeliveries
  const written = Math.floor(total * (0.62 + rand() * 0.08))
  const interview = Math.floor(written * (0.55 + rand() * 0.08))
  const waiting = Math.floor(interview * (0.42 + rand() * 0.08))
  const offer = Math.floor(waiting * (0.68 + rand() * 0.1))
  const stages = [
    { label: '投递', value: total, pct: 100, color: '#007aff', convPct: 0 },
    { label: '笔试', value: written, pct: Math.round((written / total) * 100), color: '#5e5ce6', convPct: Math.round((written / total) * 100) },
    { label: '面试', value: interview, pct: Math.round((interview / total) * 100), color: '#af52de', convPct: Math.round((interview / written) * 100) },
    { label: '待Offer', value: waiting, pct: Math.round((waiting / total) * 100), color: '#ff9500', convPct: Math.round((waiting / interview) * 100) },
    { label: 'Offer', value: offer, pct: Math.round((offer / total) * 100), color: '#34c759', convPct: Math.round((offer / waiting) * 100) },
  ]
  return stages
})

// 热门求职方向
const hotJobs = computed(() => {
  const rand = seedRand(todaySeed.value + 4)
  const data = [
    { name: '前端工程师', count: 8200 + Math.floor(rand() * 800) },
    { name: '产品经理', count: 6800 + Math.floor(rand() * 600) },
    { name: '后端工程师', count: 7400 + Math.floor(rand() * 700) },
    { name: '算法工程师', count: 5200 + Math.floor(rand() * 500) },
    { name: '数据分析师', count: 4600 + Math.floor(rand() * 400) },
  ]
  const max = Math.max(...data.map((d) => d.count))
  return data.map((d) => ({ ...d, pct: Math.round((d.count / max) * 100) }))
})

// 渠道效果对比
const channelComparison = computed(() => {
  const rand = seedRand(todaySeed.value + 5)
  const data = [
    { name: 'BOSS直聘', count: 4200 + Math.floor(rand() * 400), pct: 0, color: '#007aff' },
    { name: '内推', count: 3600 + Math.floor(rand() * 400), pct: 0, color: '#af52de' },
    { name: '官网', count: 2800 + Math.floor(rand() * 300), pct: 0, color: '#5856d6' },
    { name: '校招', count: 2200 + Math.floor(rand() * 200), pct: 0, color: '#34c759' },
    { name: '猎头', count: 1400 + Math.floor(rand() * 200), pct: 0, color: '#ff9500' },
  ]
  const max = Math.max(...data.map((d) => d.count))
  return data.map((d) => ({ ...d, pct: Math.round((d.count / max) * 100) }))
})

// 平台数据最后更新时间
const platformUpdated = computed(() => {
  return dayjs().format('YYYY-MM-DD HH:mm')
})

// ==================== 生命周期 ====================
onMounted(async () => {
  await deliveryStore.fetchDeliveries()
})

onUnmounted(() => {
  deliveryStore.clearCurrentDelivery()
})
</script>

<style scoped>
/* ==================== Pinguo 设计系统：投递看板 ==================== */
.dk-shell {
  outline: none;
  padding: 24px;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
  font-family: var(--font-sans);
}

/* ===== Quick Capture 顶部栏 ===== */
.qc-bar {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 14px 18px;
  margin-bottom: 20px;
  box-shadow: var(--shadow-sm);
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.qc-input-wrap {
  display: flex;
  align-items: center;
  gap: 10px;
}
.qc-icon {
  color: var(--primary);
  flex-shrink: 0;
  font-size: 16px;
}
.qc-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 14px;
  color: var(--foreground);
  outline: none;
  padding: 4px 0;
  font-family: inherit;
}
.qc-input::placeholder {
  color: var(--muted-foreground);
}
.qc-kbd {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 600;
  color: var(--muted-foreground);
  background: var(--background-200);
  border-radius: 6px;
  padding: 3px 8px;
}
.qc-hint {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.qc-parsed-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  color: var(--primary);
  background: var(--brand-50);
  padding: 3px 10px;
  border-radius: 9999px;
}

/* ===== 统计卡片 ===== */
.stat-section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--muted-foreground);
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.stat-section-title::before {
  content: '';
  width: 3px;
  height: 14px;
  background: var(--primary);
  border-radius: 2px;
}
.stat-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 22px;
}
.stat-card {
  position: relative;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  box-shadow: var(--shadow-sm);
  overflow: hidden;
  transition: box-shadow 0.2s ease, transform 0.2s ease;
}
.stat-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}
.stat-icon-wrap {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 20px;
}
.stat-body {
  flex: 1;
  min-width: 0;
}
.stat-number {
  font-family: var(--font-mono);
  font-size: 26px;
  font-weight: 700;
  line-height: 1;
  color: var(--foreground);
  letter-spacing: -0.02em;
}
.stat-label {
  font-size: 12px;
  color: var(--muted-foreground);
  margin-top: 5px;
}
.stat-change {
  font-size: 11px;
  color: var(--muted-foreground);
  margin-top: 5px;
  font-family: var(--font-mono);
}
.stat-change.up {
  color: var(--success);
}
.stat-change.down {
  color: var(--destructive);
}
.stat-spark {
  position: absolute;
  right: 16px;
  bottom: 16px;
  width: 80px;
  height: 24px;
  opacity: 0.6;
  pointer-events: none;
}

/* ===== 视图切换 + 筛选器 ===== */
.section-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
  flex-wrap: wrap;
  gap: 12px;
}
.view-toggle-group {
  display: inline-flex;
  background: var(--background-200);
  border-radius: 10px;
  padding: 3px;
  gap: 2px;
}
.view-toggle-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 6px 14px;
  border-radius: 7px;
  font-size: 13px;
  font-weight: 600;
  color: var(--muted-foreground);
  transition: background 0.15s ease, color 0.15s ease, box-shadow 0.15s ease;
  font-family: inherit;
}
.view-toggle-btn[data-active='true'] {
  background: var(--card);
  color: var(--foreground);
  box-shadow: var(--shadow-sm);
}
.search-wrap {
  position: relative;
  display: inline-flex;
  align-items: center;
}
.search-icon {
  position: absolute;
  left: 12px;
  color: var(--muted-foreground);
  pointer-events: none;
  font-size: 14px;
}
.search-input {
  width: 220px;
  height: 36px;
  border: 1px solid var(--border);
  border-radius: 10px;
  background: var(--card);
  padding: 0 12px 0 34px;
  font-size: 13px;
  color: var(--foreground);
  outline: none;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
  font-family: inherit;
}
.search-input::placeholder {
  color: var(--muted-foreground);
}
.search-input:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--brand-50);
}
.filter-cluster {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.filter-select {
  height: 36px;
  border: 1px solid var(--border);
  border-radius: 10px;
  background: var(--card);
  padding: 0 32px 0 12px;
  font-size: 13px;
  color: var(--foreground);
  outline: none;
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='14' height='14' viewBox='0 0 24 24' fill='none' stroke='%238e8e93' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'><polyline points='6 9 12 15 18 9'></polyline></svg>");
  background-repeat: no-repeat;
  background-position: right 10px center;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
  font-family: inherit;
}
.filter-select:focus {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--brand-50);
}
.btn-primary-capsule {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 36px;
  padding: 0 16px;
  border-radius: 9999px;
  background: var(--primary);
  color: var(--primary-foreground);
  font-size: 13px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: background-color 0.15s ease, box-shadow 0.15s ease, transform 0.15s ease;
  font-family: inherit;
}
.btn-primary-capsule:hover {
  background: var(--brand-600);
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}
.btn-secondary-capsule {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 36px;
  padding: 0 16px;
  border-radius: 9999px;
  background: var(--background-200);
  color: var(--foreground);
  font-size: 13px;
  font-weight: 600;
  border: 1px solid var(--border);
  cursor: pointer;
  transition: background-color 0.15s ease;
  font-family: inherit;
}
.btn-secondary-capsule:hover {
  background: var(--background-300);
}

/* ===== 主体布局 ===== */
.content-flex {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}
.content-flex.platform-mode {
  flex-direction: column;
}
.content-main {
  flex: 1;
  min-width: 0;
}

/* ===== 表格 ===== */
.kb-table-card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 16px;
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}
.kb-table-scroll {
  overflow-x: auto;
}
.kb-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  min-width: 980px;
}
.kb-table thead th {
  background: var(--background-200);
  color: var(--muted-foreground);
  font-weight: 600;
  text-align: left;
  padding: 12px 14px;
  font-size: 12px;
  white-space: nowrap;
  border-bottom: 1px solid var(--border);
  letter-spacing: 0.02em;
}
.kb-table tbody td {
  padding: 14px;
  border-bottom: 1px solid var(--border);
  vertical-align: middle;
  color: var(--foreground);
}
.kb-table tbody tr:last-child td {
  border-bottom: none;
}
.kb-table tbody tr {
  transition: background 0.15s ease;
  cursor: pointer;
}
.kb-table tbody tr:hover {
  background: var(--background-100);
}
.kb-table tbody tr.selected {
  background: var(--brand-50);
}
.kb-table tbody tr.selected td:first-child {
  box-shadow: inset 3px 0 0 var(--primary);
}
.kb-table .mono {
  font-family: var(--font-mono);
  color: var(--muted-foreground);
  font-size: 12px;
}
.kb-company-cell {
  display: flex;
  align-items: center;
  gap: 4px;
}
.kb-company {
  font-weight: 700;
  color: var(--foreground);
}
.kb-remark-dot {
  color: var(--muted-foreground);
}
.kb-feedback-cell {
  font-size: 12px;
  color: var(--muted-foreground);
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.kb-dash {
  color: var(--background-400);
}

/* 渠道标签 */
.kb-channel {
  display: inline-block;
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 9999px;
  white-space: nowrap;
}
.ch-boss {
  background: var(--brand-50);
  color: var(--primary);
}
.ch-referral {
  background: rgba(175, 82, 222, 0.14);
  color: #8e3db5;
}
.ch-official {
  background: rgba(88, 86, 214, 0.14);
  color: #5856d6;
}
.ch-campus {
  background: var(--state-success-surface);
  color: var(--success);
}
.ch-headhunt {
  background: rgba(255, 149, 0, 0.14);
  color: #c46900;
}
.ch-other {
  background: var(--background-200);
  color: var(--muted-foreground);
}

/* 下次面试 */
.kb-next-interview {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 600;
  color: var(--primary);
  background: var(--brand-50);
  padding: 3px 8px;
  border-radius: 9999px;
  white-space: nowrap;
}

/* 状态标签 */
.kb-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 3px 10px;
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}
.kb-status .status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}
.kb-status.gray {
  background: rgba(174, 174, 178, 0.18);
  color: var(--background-600);
}
.kb-status.gray .status-dot {
  background: var(--background-600);
}
.kb-status.blue {
  background: var(--brand-50);
  color: var(--primary);
}
.kb-status.blue .status-dot {
  background: var(--primary);
}
.kb-status.orange {
  background: rgba(255, 149, 0, 0.14);
  color: #c46900;
}
.kb-status.orange .status-dot {
  background: #ff9500;
}
.kb-status.purple {
  background: rgba(175, 82, 222, 0.14);
  color: #8e3db5;
}
.kb-status.purple .status-dot {
  background: #af52de;
}
.kb-status.green {
  background: var(--state-success-surface);
  color: var(--success);
}
.kb-status.green .status-dot {
  background: var(--success);
}
.kb-status.red {
  background: var(--state-error-surface);
  color: var(--destructive);
}
.kb-status.red .status-dot {
  background: var(--destructive);
}

/* 进度点 */
.kb-dots {
  display: inline-flex;
  align-items: center;
  gap: 5px;
}
.kb-dots .dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--background-400);
}
.kb-dots .dot.done {
  background: var(--primary);
}
.kb-dots .dot.current {
  background: transparent;
  box-shadow: 0 0 0 2px var(--primary) inset;
}
.kb-dots .dot.fail {
  background: var(--destructive);
}

/* 优先级 */
.kb-priority {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
  padding: 3px 9px;
  border-radius: 9999px;
}
.kb-priority.high {
  background: rgba(255, 59, 48, 0.12);
  color: var(--destructive);
}
.kb-priority.mid {
  background: rgba(255, 149, 0, 0.12);
  color: #c46900;
}
.kb-priority.low {
  background: var(--background-200);
  color: var(--muted-foreground);
}

/* 操作按钮 */
.kb-action {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  border: none;
  background: transparent;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--muted-foreground);
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease;
  font-size: 14px;
}
.kb-action:hover {
  background: var(--background-200);
  color: var(--foreground);
}
.kb-action-sm {
  width: 22px;
  height: 22px;
  font-size: 12px;
  border-radius: 6px;
}

/* ===== 移动端卡片 ===== */
.kb-cards-mobile {
  display: none;
}
.kb-mobile-card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 14px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: box-shadow 0.15s ease;
}
.kb-mobile-card:hover {
  box-shadow: var(--shadow-sm);
}
.kb-mobile-card.selected {
  border-color: var(--primary);
  box-shadow: 0 0 0 2px var(--brand-50);
}
.kb-mobile-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}
.kb-mobile-pos {
  font-size: 13px;
  color: var(--foreground);
  margin-bottom: 8px;
}
.kb-mobile-meta {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
  font-size: 11px;
}
.kb-mobile-foot {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.kb-mobile-fb {
  font-size: 11px;
  color: var(--muted-foreground);
}

/* ===== 平台数据视图 ===== */
.platform-view {
  display: flex;
  flex-direction: column;
  gap: 18px;
}
.platform-hero {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
}
.platform-hero-card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 18px 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  box-shadow: var(--shadow-sm);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}
.platform-hero-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}
.platform-hero-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 22px;
}
.platform-hero-value {
  font-family: var(--font-mono);
  font-size: 26px;
  font-weight: 800;
  line-height: 1.1;
  color: var(--foreground);
  letter-spacing: -0.02em;
}
.platform-hero-label {
  font-size: 12px;
  color: var(--muted-foreground);
  margin-top: 3px;
}
.platform-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
.platform-card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 20px;
  box-shadow: var(--shadow-sm);
}
.platform-card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
}
.platform-card-title {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 700;
  color: var(--foreground);
}
.platform-card-sub {
  font-size: 11px;
  color: var(--muted-foreground);
  font-family: var(--font-mono);
}

/* 饼图 */
.pie-wrap {
  display: flex;
  gap: 20px;
  align-items: center;
  flex-wrap: wrap;
}
.pie-svg {
  width: 180px;
  height: 180px;
  flex-shrink: 0;
}
.pie-center-num {
  font-family: var(--font-mono);
  font-size: 22px;
  font-weight: 800;
  fill: var(--foreground);
  letter-spacing: -0.02em;
}
.pie-center-label {
  font-size: 11px;
  fill: var(--muted-foreground);
}
.pie-legend {
  flex: 1;
  min-width: 160px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.pie-legend-item {
  display: grid;
  grid-template-columns: 12px 1fr auto auto;
  gap: 8px;
  align-items: center;
  font-size: 12px;
  padding: 6px 10px;
  background: var(--background-100);
  border-radius: 10px;
}
.pie-legend-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}
.pie-legend-name {
  color: var(--foreground);
  font-weight: 600;
}
.pie-legend-value {
  font-family: var(--font-mono);
  color: var(--muted-foreground);
  font-size: 11px;
}
.pie-legend-pct {
  font-family: var(--font-mono);
  color: var(--foreground);
  font-weight: 700;
  font-size: 12px;
}

/* 漏斗 */
.funnel-platform {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.funnel-platform-row {
  display: grid;
  grid-template-columns: 80px 1fr auto auto auto;
  gap: 12px;
  align-items: center;
}
.funnel-stage-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--foreground);
}
.funnel-stage-bar-wrap {
  height: 24px;
  background: var(--background-100);
  border-radius: 9999px;
  overflow: hidden;
}
.funnel-stage-bar {
  height: 100%;
  border-radius: 9999px;
  transition: width 0.4s ease;
}
.funnel-stage-value {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  color: var(--foreground);
  min-width: 60px;
  text-align: right;
}
.funnel-stage-pct {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--muted-foreground);
  min-width: 48px;
  text-align: right;
}
.funnel-stage-conv {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  color: var(--success);
  display: inline-flex;
  align-items: center;
  gap: 3px;
  min-width: 56px;
  background: var(--state-success-surface);
  padding: 3px 8px;
  border-radius: 9999px;
}

/* 热门岗位 */
.hot-job-row {
  display: grid;
  grid-template-columns: 24px 1fr 1fr auto;
  gap: 10px;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--border);
}
.hot-job-row:last-child {
  border-bottom: none;
}
.hot-job-rank {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 800;
  font-family: var(--font-mono);
  background: var(--background-200);
  color: var(--muted-foreground);
}
.hot-job-rank.rank-1 {
  background: linear-gradient(135deg, #ff9500, #af52de);
  color: #fff;
}
.hot-job-rank.rank-2 {
  background: linear-gradient(135deg, var(--background-400), var(--background-600));
  color: #fff;
}
.hot-job-rank.rank-3 {
  background: linear-gradient(135deg, #ff9500, #5856d6);
  color: #fff;
}
.hot-job-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--foreground);
}
.hot-job-bar-wrap {
  height: 6px;
  background: var(--background-100);
  border-radius: 9999px;
  overflow: hidden;
}
.hot-job-bar {
  height: 100%;
  background: linear-gradient(90deg, var(--primary), var(--brand-400));
  border-radius: 9999px;
  transition: width 0.4s ease;
}
.hot-job-count {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  color: var(--muted-foreground);
  min-width: 56px;
  text-align: right;
}

/* 渠道对比 */
.channel-cmp-row {
  display: grid;
  grid-template-columns: 80px 1fr auto auto;
  gap: 10px;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--border);
}
.channel-cmp-row:last-child {
  border-bottom: none;
}
.channel-cmp-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--foreground);
}
.channel-cmp-bar-wrap {
  height: 6px;
  background: var(--background-100);
  border-radius: 9999px;
  overflow: hidden;
}
.channel-cmp-bar {
  height: 100%;
  border-radius: 9999px;
  transition: width 0.4s ease;
}
.channel-cmp-pct {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  color: var(--foreground);
  min-width: 40px;
  text-align: right;
}
.channel-cmp-count {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--muted-foreground);
  min-width: 60px;
  text-align: right;
}

.platform-footer {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px 16px;
  background: var(--background-100);
  border-radius: 10px;
  font-size: 11px;
  color: var(--muted-foreground);
}
.platform-update {
  margin-left: auto;
  font-family: var(--font-mono);
}

/* ===== 详情面板 ===== */
.detail-panel {
  width: 380px;
  flex-shrink: 0;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 16px;
  box-shadow: var(--shadow-sm);
  padding: 22px;
  position: sticky;
  top: 24px;
  align-self: flex-start;
  max-height: calc(100vh - 48px);
  overflow-y: auto;
}
.detail-panel::-webkit-scrollbar {
  width: 6px;
}
.detail-panel::-webkit-scrollbar-track {
  background: transparent;
}
.detail-panel::-webkit-scrollbar-thumb {
  background: var(--background-300);
  border-radius: 9999px;
}
.detail-panel.detail-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 40px 22px;
}
.detail-empty-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: var(--brand-50);
  color: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin-bottom: 16px;
}
.detail-empty-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--foreground);
  margin-bottom: 6px;
}
.detail-empty-desc {
  font-size: 12px;
  color: var(--muted-foreground);
  line-height: 1.5;
  max-width: 240px;
}
.detail-empty-shortcuts {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 11px;
  color: var(--muted-foreground);
}
.shortcut-row {
  display: flex;
  align-items: center;
  gap: 6px;
  justify-content: center;
}
.shortcut-row kbd {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 600;
  background: var(--background-200);
  border-radius: 4px;
  padding: 2px 6px;
  color: var(--foreground);
}

.dp-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}
.dp-header-left {
  min-width: 0;
  flex: 1;
}
.dp-company-name {
  font-size: 18px;
  font-weight: 800;
  color: var(--foreground);
  line-height: 1.3;
  letter-spacing: -0.01em;
}
.dp-company-sub {
  font-size: 13px;
  color: var(--muted-foreground);
  margin-top: 3px;
}
.dp-close {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: var(--muted-foreground);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: background 0.15s ease, color 0.15s ease;
  font-size: 14px;
}
.dp-close:hover {
  background: var(--background-200);
  color: var(--foreground);
}
.dp-tags {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  flex-wrap: wrap;
}
.dp-channel-tag {
  font-size: 11px;
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 9999px;
  background: var(--background-200);
  color: var(--secondary-foreground);
}
.dp-meta-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 12px;
  padding: 10px 12px;
  background: var(--background-100);
  border-radius: 10px;
}
.dp-meta-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: var(--muted-foreground);
  font-family: var(--font-mono);
}
.dp-offer-card {
  margin-top: 14px;
  padding: 14px;
  background: linear-gradient(
    135deg,
    var(--state-success-surface) 0%,
    var(--card) 70%
  );
  border: 1px solid var(--border);
  border-radius: 10px;
}
.dp-offer-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 700;
  color: var(--success);
}
.dp-offer-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-top: 10px;
}
.dp-offer-item {
  background: var(--card);
  border-radius: 10px;
  padding: 8px 10px;
}
.dp-offer-label {
  font-size: 10px;
  color: var(--muted-foreground);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.dp-offer-value {
  font-family: var(--font-mono);
  font-size: 14px;
  font-weight: 700;
  color: var(--foreground);
  margin-top: 2px;
}
.dp-section-title {
  font-size: 11px;
  font-weight: 700;
  color: var(--muted-foreground);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin: 20px 0 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.dp-add-btn {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  border: none;
  background: transparent;
  color: var(--primary);
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  padding: 3px 7px;
  border-radius: 6px;
  text-transform: none;
  letter-spacing: 0;
  font-family: inherit;
}
.dp-add-btn:hover {
  background: var(--brand-50);
}

/* 时间线 */
.tl-list {
  display: flex;
  flex-direction: column;
}
.tl-item {
  display: flex;
  gap: 10px;
  padding-bottom: 14px;
  position: relative;
}
.tl-item::before {
  content: '';
  position: absolute;
  left: 5px;
  top: 14px;
  bottom: 0;
  width: 1px;
  background: var(--border);
}
.tl-item:last-child::before {
  display: none;
}
.tl-item.line-dashed::before {
  background: transparent;
  border-left: 1px dashed var(--border);
}
.tl-item.line-none::before {
  display: none;
}
.tl-dot {
  width: 11px;
  height: 11px;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 3px;
  position: relative;
  z-index: 1;
  background: var(--card);
  box-shadow: 0 0 0 2px var(--background-400) inset;
}
.tl-dot.gray {
  background: var(--background-400);
  box-shadow: none;
}
.tl-dot.green {
  background: var(--success);
  box-shadow: none;
}
.tl-dot.blue {
  background: var(--primary);
  box-shadow: none;
}
.tl-dot.red {
  background: var(--destructive);
  box-shadow: none;
}
.tl-content {
  flex: 1;
  min-width: 0;
}
.round-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.tl-title {
  font-size: 13px;
  font-weight: 700;
  color: var(--foreground);
}
.tl-time {
  font-size: 11px;
  color: var(--muted-foreground);
  margin-top: 2px;
  font-family: var(--font-mono);
}
.tl-meta {
  font-size: 11px;
  color: var(--muted-foreground);
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
}
.tl-tag {
  display: inline-block;
  margin-top: 4px;
  font-size: 11px;
  color: var(--primary);
  background: var(--brand-50);
  padding: 2px 8px;
  border-radius: 9999px;
}
.tl-action-row {
  margin-top: 6px;
}
.result-pass {
  color: var(--success);
  font-weight: 600;
}
.result-fail {
  color: var(--destructive);
  font-weight: 600;
}
.result-pending {
  color: var(--muted-foreground);
}
.dp-link-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 6px;
  border: none;
  background: transparent;
  color: var(--primary);
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  font-family: inherit;
}
.dp-link-btn:hover {
  background: var(--brand-50);
}

/* 反馈 */
.fb-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.fb-bubble {
  background: var(--background-100);
  border-radius: 10px;
  padding: 10px 12px;
}
.fb-text {
  font-size: 12px;
  color: var(--foreground);
  line-height: 1.5;
}
.fb-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 6px;
  font-size: 10px;
  color: var(--muted-foreground);
  font-family: var(--font-mono);
}
.fb-meta span {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.fb-del {
  border: none;
  background: transparent;
  color: var(--muted-foreground);
  cursor: pointer;
  padding: 2px;
  border-radius: 4px;
  font-size: 12px;
}
.fb-del:hover {
  background: var(--state-error-surface);
  color: var(--destructive);
}
.fb-next {
  margin-top: 6px;
  font-size: 11px;
  color: var(--primary);
  display: inline-flex;
  align-items: center;
  gap: 3px;
}
.empty-hint {
  font-size: 12px;
  color: var(--muted-foreground);
  padding: 14px;
  text-align: center;
  background: var(--background-100);
  border-radius: 10px;
}

/* 策略卡片 */
.strategy-card {
  margin-top: 18px;
  padding: 14px;
  background: var(--background-100);
  border-radius: 12px;
}
.strategy-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 700;
  color: var(--foreground);
}
.strategy-sub {
  font-size: 11px;
  font-weight: 600;
  color: var(--muted-foreground);
  margin: 12px 0 8px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.funnel-svg-wrap {
  display: flex;
  gap: 12px;
  align-items: center;
}
.funnel-svg {
  width: 120px;
  height: 80px;
  flex-shrink: 0;
}
.funnel-labels {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 3px;
}
.funnel-label-item {
  display: flex;
  justify-content: space-between;
  font-size: 10px;
}
.funnel-label-name {
  color: var(--muted-foreground);
}
.funnel-label-val {
  font-family: var(--font-mono);
  color: var(--foreground);
  font-weight: 600;
}
.funnel-label-pct {
  color: var(--muted-foreground);
}
.bar-row {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 8px;
}
.bar-label {
  font-size: 11px;
  color: var(--foreground);
  display: flex;
  justify-content: space-between;
}
.bar-label .mono {
  font-family: var(--font-mono);
  color: var(--muted-foreground);
}
.bar-track {
  height: 6px;
  background: var(--background-200);
  border-radius: 9999px;
  overflow: hidden;
}
.bar-fill {
  height: 100%;
  border-radius: 9999px;
  background: var(--primary);
  transition: width 0.4s ease;
}
.bar-fill.boss {
  background: var(--primary);
}
.bar-fill.web {
  background: #5856d6;
}
.bar-fill.ref {
  background: #af52de;
}
.bar-fill.you {
  background: var(--success);
}
.bar-fill.avg {
  background: var(--background-500);
}
.cycle-stat {
  display: flex;
  align-items: baseline;
  gap: 6px;
}
.cycle-num {
  font-family: var(--font-mono);
  font-size: 28px;
  font-weight: 800;
  color: var(--primary);
  letter-spacing: -0.02em;
}
.cycle-unit {
  font-size: 11px;
  color: var(--muted-foreground);
}

/* ===== 空状态 ===== */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 80px 24px;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 16px;
}
.empty-state-full {
  min-height: 400px;
}
.empty-illustration {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--brand-50);
  color: var(--primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  margin-bottom: 16px;
}
.empty-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--foreground);
  margin-bottom: 6px;
}
.empty-desc {
  font-size: 14px;
  color: var(--muted-foreground);
  margin-bottom: 20px;
  max-width: 360px;
}

/* ===== 视图切换动画 ===== */
.view-fade-enter-active,
.view-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}
.view-fade-enter-from {
  opacity: 0;
  transform: translateY(8px);
}
.view-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* 详情面板滑入动画 */
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: opacity 0.25s ease, transform 0.25s cubic-bezier(0.32, 0.72, 0, 1);
}
.slide-fade-enter-from {
  opacity: 0;
  transform: translateX(20px);
}
.slide-fade-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

/* ===== 响应式 ===== */
@media (max-width: 1024px) {
  .dk-shell {
    padding: 16px;
  }
  .stat-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .platform-hero {
    grid-template-columns: repeat(2, 1fr);
  }
  .platform-grid-2 {
    grid-template-columns: 1fr;
  }
  .detail-panel {
    display: none;
  }
}

@media (max-width: 768px) {
  .dk-shell {
    padding: 12px;
  }
  .kb-table-desktop {
    display: none;
  }
  .kb-cards-mobile {
    display: block;
  }
  .section-row {
    flex-direction: column;
    align-items: stretch;
  }
  .search-wrap {
    width: 100%;
  }
  .search-input {
    width: 100%;
  }
  .filter-cluster {
    width: 100%;
    justify-content: space-between;
  }
  .filter-select {
    flex: 1;
    min-width: 0;
  }
  .btn-primary-capsule {
    flex-shrink: 0;
  }
}

@media (max-width: 640px) {
  .stat-grid {
    grid-template-columns: 1fr;
  }
  .platform-hero {
    grid-template-columns: 1fr;
  }
  .funnel-platform-row {
    grid-template-columns: 70px 1fr auto auto;
  }
  .funnel-stage-conv {
    display: none;
  }
  .pie-wrap {
    flex-direction: column;
    align-items: stretch;
  }
  .pie-svg {
    margin: 0 auto;
  }
}

/* 无障碍：尊重减少动效偏好 */
@media (prefers-reduced-motion: reduce) {
  .stat-card,
  .platform-hero-card,
  .btn-primary-capsule,
  .btn-secondary-capsule,
  .kb-action,
  .view-toggle-btn {
    transition: none;
  }
  .stat-card:hover,
  .platform-hero-card:hover {
    transform: none;
  }
}
</style>