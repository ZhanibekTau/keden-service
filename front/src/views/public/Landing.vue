<template>
  <div class="landing">
    <!-- Header -->
    <header class="header">
      <div class="container header-content">
        <div class="logo-group">
          <h1 class="logo">KEDEN</h1>
          <span class="logo-divider"></span>
          <span class="logo-sub">AI-powered</span>
        </div>
        <nav class="nav-links">
          <a href="#how-it-works" class="nav-link">Как это работает</a>
          <a href="#benefits" class="nav-link">Преимущества</a>
          <a href="#pricing" class="nav-link">Стоимость</a>
          <router-link to="/login"><el-button text class="nav-btn">Войти</el-button></router-link>
          <el-button type="primary" class="register-btn" @click="showRegisterModal = true">Зарегистрироваться</el-button>
        </nav>
      </div>
    </header>

    <!-- Hero -->
    <section class="hero">
      <div class="container hero-inner">
        <div class="hero-badge">AI-powered платформа</div>
        <h1>Платформа для<br/>таможенных брокеров</h1>
        <p>Автоматизация обработки таможенных документов с использованием AI и выдачей структурированных Excel-файлов</p>
        <div class="hero-actions">
          <el-button type="primary" size="large" class="hero-btn" @click="showRegisterModal = true">Начать работу</el-button>
          <a href="#how-it-works">
            <el-button size="large" class="hero-btn-outline">Узнать больше</el-button>
          </a>
        </div>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-number">500+</span>
            <span class="stat-label">Документов обработано</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">50+</span>
            <span class="stat-label">Компаний</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">99%</span>
            <span class="stat-label">Точность AI</span>
          </div>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section id="how-it-works" class="section">
      <div class="container">
        <h2 class="section-title">Как это работает</h2>
        <p class="section-subtitle">Три простых шага для обработки документов</p>
        <el-row :gutter="40" justify="center">
          <el-col :span="8">
            <div class="step-card">
              <div class="step-icon" style="background: rgba(38, 166, 91, 0.1)">
                <el-icon :size="32" color="#26A65B"><Upload /></el-icon>
              </div>
              <div class="step-number">01</div>
              <h3>Загрузите PDF</h3>
              <p>Загрузите пакет таможенных документов в формате PDF</p>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="step-card">
              <div class="step-icon" style="background: rgba(38, 166, 91, 0.1)">
                <el-icon :size="32" color="#26A65B"><Cpu /></el-icon>
              </div>
              <div class="step-number">02</div>
              <h3>AI обработка</h3>
              <p>Искусственный интеллект анализирует и извлекает данные из документов</p>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="step-card">
              <div class="step-icon" style="background: rgba(38, 166, 91, 0.1)">
                <el-icon :size="32" color="#26A65B"><Download /></el-icon>
              </div>
              <div class="step-number">03</div>
              <h3>Скачайте Excel</h3>
              <p>Получите структурированный Excel-файл с результатами обработки</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </section>

    <!-- Benefits -->
    <section id="benefits" class="section section-white">
      <div class="container">
        <h2 class="section-title">Преимущества</h2>
        <p class="section-subtitle">Почему компании выбирают KEDEN</p>
        <el-row :gutter="30">
          <el-col :span="6" v-for="item in benefits" :key="item.title">
            <div class="benefit-card">
              <div class="benefit-icon" :style="{ background: item.bg }">
                <el-icon :size="28" :color="item.color"><component :is="item.icon" /></el-icon>
              </div>
              <h4>{{ item.title }}</h4>
              <p>{{ item.desc }}</p>
            </div>
          </el-col>
        </el-row>
      </div>
    </section>

    <!-- Pricing -->
    <section id="pricing" class="section">
      <div class="container" style="text-align: center">
        <h2 class="section-title">Стоимость</h2>
        <p class="section-subtitle">Простая и прозрачная тарификация</p>
        <div class="pricing-card">
          <div class="pricing-badge">Единый тариф</div>
          <h2 class="pricing-amount">12 990 &#8376;</h2>
          <p class="pricing-period">в месяц</p>
          <ul class="pricing-features">
            <li><el-icon color="#26A65B"><Check /></el-icon> Неограниченная загрузка документов</li>
            <li><el-icon color="#26A65B"><Check /></el-icon> AI-обработка</li>
            <li><el-icon color="#26A65B"><Check /></el-icon> Генерация Excel-файлов</li>
            <li><el-icon color="#26A65B"><Check /></el-icon> Техническая поддержка</li>
          </ul>
          <el-button type="primary" size="large" style="width: 100%" @click="showRegisterModal = true">Подать заявку</el-button>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="footer">
      <div class="container">
        <el-row :gutter="40">
          <el-col :span="8">
            <div class="footer-section">
              <h3 class="footer-logo">KEDEN</h3>
              <p class="footer-desc">Платформа для автоматизации обработки таможенных документов с применением искусственного интеллекта.</p>
            </div>
          </el-col>
          <el-col :span="5">
            <div class="footer-section">
              <h4>Навигация</h4>
              <ul class="footer-links">
                <li><a href="#how-it-works">Как это работает</a></li>
                <li><a href="#benefits">Преимущества</a></li>
                <li><a href="#pricing">Стоимость</a></li>
                <li><router-link to="/login">Войти</router-link></li>
              </ul>
            </div>
          </el-col>
          <el-col :span="5">
            <div class="footer-section">
              <h4>Правовая информация</h4>
              <ul class="footer-links">
                <li><a href="#">Политика конфиденциальности</a></li>
                <li><a href="#">Условия использования</a></li>
                <li><a href="#">Публичная оферта</a></li>
              </ul>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="footer-section">
              <h4>Контакты</h4>
              <ul class="footer-contacts">
                <li>
                  <el-icon><Phone /></el-icon>
                  <span>+7 (727) 123-45-67</span>
                </li>
                <li>
                  <el-icon><Message /></el-icon>
                  <span>info@keden.kz</span>
                </li>
                <li>
                  <el-icon><Location /></el-icon>
                  <span>г. Алматы, ул. Абая 52, офис 301</span>
                </li>
                <li>
                  <el-icon><Clock /></el-icon>
                  <span>Пн-Пт: 09:00 - 18:00</span>
                </li>
              </ul>
            </div>
          </el-col>
        </el-row>
        <div class="footer-bottom">
          <p>&copy; 2026 KEDEN. Все права защищены.</p>
        </div>
      </div>
    </footer>

    <!-- Registration Type Modal -->
    <el-dialog v-model="showRegisterModal" title="Выберите тип регистрации" width="500px" :close-on-click-modal="true" align-center>
      <div class="register-type-options">
        <div class="register-type-card" @click="goToRegister('individual')">
          <el-icon :size="48" color="#26A65B"><User /></el-icon>
          <h3>Физическое лицо</h3>
          <p>Регистрация как индивидуальный брокер без привязки к компании</p>
        </div>
        <div class="register-type-card" @click="goToRegister('company')">
          <el-icon :size="48" color="#26A65B"><OfficeBuilding /></el-icon>
          <h3>Компания</h3>
          <p>Регистрация как юридическое лицо с данными компании</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Upload, Download, Cpu, Timer, Lock, TrendCharts, DocumentChecked, Check, Phone, Message, Location, Clock, User, OfficeBuilding } from '@element-plus/icons-vue'

const router = useRouter()
const showRegisterModal = ref(false)

const benefits = [
  { title: 'Экономия времени', desc: 'Обработка документов за минуты вместо часов', icon: Timer, color: '#26A65B', bg: 'rgba(38, 166, 91, 0.1)' },
  { title: 'Точность', desc: 'AI обеспечивает высокую точность извлечения данных', icon: DocumentChecked, color: '#1B9C85', bg: 'rgba(27, 156, 133, 0.1)' },
  { title: 'Безопасность', desc: 'Все данные защищены и хранятся на защищенных серверах', icon: Lock, color: '#2D8B55', bg: 'rgba(45, 139, 85, 0.1)' },
  { title: 'Аналитика', desc: 'Отслеживание статусов обработки в реальном времени', icon: TrendCharts, color: '#0EA47A', bg: 'rgba(14, 164, 122, 0.1)' }
]

function goToRegister(type: string) {
  showRegisterModal.value = false
  router.push({ path: '/register', query: { type } })
}
</script>

<style scoped>
.landing { background: #fff; }
.container { max-width: 1200px; margin: 0 auto; padding: 0 20px; }

/* Header */
.header {
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  padding: 14px 0;
  position: sticky;
  top: 0;
  z-index: 10;
}
.header-content { display: flex; align-items: center; justify-content: space-between; }
.logo-group { display: flex; align-items: center; gap: 12px; }
.logo { font-size: 24px; color: #26A65B; margin: 0; font-weight: 800; letter-spacing: 2px; }
.logo-divider { width: 1px; height: 20px; background: #dcdfe6; }
.logo-sub { color: #909399; font-size: 13px; font-weight: 400; }
.nav-links { display: flex; align-items: center; gap: 8px; }
.nav-link { color: #606266; text-decoration: none; font-size: 14px; padding: 6px 12px; border-radius: 6px; transition: color 0.2s; }
.nav-link:hover { color: #26A65B; }
.nav-btn { color: #606266 !important; }
.register-btn { background: #26A65B !important; border-color: #26A65B !important; }
.register-btn:hover { background: #1f8f4e !important; border-color: #1f8f4e !important; }

/* Hero */
.hero {
  background: linear-gradient(160deg, #f0faf4 0%, #fff 40%, #fff 60%, #f0faf4 100%);
  text-align: center;
  padding: 100px 0 80px;
  position: relative;
}
.hero-inner { max-width: 800px; margin: 0 auto; }
.hero-badge {
  display: inline-block;
  background: rgba(38, 166, 91, 0.1);
  color: #26A65B;
  font-size: 13px;
  font-weight: 600;
  padding: 6px 16px;
  border-radius: 20px;
  margin-bottom: 24px;
}
.hero h1 { font-size: 48px; margin-bottom: 20px; color: #1a1a2e; font-weight: 800; line-height: 1.2; }
.hero p { font-size: 18px; margin-bottom: 36px; color: #606266; line-height: 1.6; }
.hero-actions { display: flex; justify-content: center; gap: 16px; margin-bottom: 60px; }
.hero-btn { background: #26A65B !important; border-color: #26A65B !important; padding: 12px 32px !important; font-size: 16px !important; }
.hero-btn:hover { background: #1f8f4e !important; border-color: #1f8f4e !important; }
.hero-btn-outline { border-color: #c0c4cc !important; color: #606266 !important; padding: 12px 32px !important; font-size: 16px !important; }
.hero-btn-outline:hover { border-color: #26A65B !important; color: #26A65B !important; }
.hero-stats {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 40px;
  background: #fff;
  border-radius: 16px;
  padding: 32px 48px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
}
.stat-item { text-align: center; }
.stat-number { display: block; font-size: 28px; font-weight: 800; color: #26A65B; }
.stat-label { display: block; font-size: 13px; color: #909399; margin-top: 4px; }
.stat-divider { width: 1px; height: 40px; background: #ebeef5; }

/* Sections */
.section { padding: 80px 0; background: #fafbfc; }
.section-white { background: #fff; }
.section-title { text-align: center; font-size: 32px; margin-bottom: 8px; color: #1a1a2e; font-weight: 700; }
.section-subtitle { text-align: center; color: #909399; font-size: 16px; margin-bottom: 48px; }

/* Step cards */
.step-card {
  text-align: center;
  padding: 32px 24px;
  background: #fff;
  border-radius: 16px;
  transition: transform 0.2s, box-shadow 0.2s;
  border: 1px solid #f0f0f0;
}
.step-card:hover { transform: translateY(-4px); box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08); }
.step-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  border-radius: 16px;
  margin-bottom: 8px;
}
.step-number { font-size: 13px; color: #26A65B; font-weight: 700; margin-bottom: 8px; }
.step-card h3 { margin: 8px 0; color: #1a1a2e; font-size: 18px; }
.step-card p { color: #909399; font-size: 14px; line-height: 1.6; }

/* Benefit cards */
.benefit-card {
  text-align: center;
  padding: 32px 20px;
  background: #fafbfc;
  border-radius: 16px;
  transition: transform 0.2s, box-shadow 0.2s;
  border: 1px solid #f0f0f0;
}
.benefit-card:hover { transform: translateY(-4px); box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08); }
.benefit-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 14px;
  margin-bottom: 12px;
}
.benefit-card h4 { margin: 8px 0; color: #1a1a2e; font-size: 16px; }
.benefit-card p { color: #909399; font-size: 13px; line-height: 1.5; }

/* Pricing */
.pricing-card {
  max-width: 420px;
  margin: 0 auto;
  background: #fff;
  border-radius: 20px;
  padding: 40px 36px;
  border: 2px solid #26A65B;
  box-shadow: 0 8px 30px rgba(38, 166, 91, 0.1);
}
.pricing-badge {
  display: inline-block;
  background: rgba(38, 166, 91, 0.1);
  color: #26A65B;
  font-size: 13px;
  font-weight: 600;
  padding: 4px 14px;
  border-radius: 12px;
  margin-bottom: 16px;
}
.pricing-amount { font-size: 40px; color: #26A65B; margin: 0; font-weight: 800; }
.pricing-period { color: #909399; margin-bottom: 24px; }
.pricing-features {
  text-align: left;
  list-style: none;
  padding: 0;
  margin: 0 0 28px;
}
.pricing-features li {
  padding: 10px 0;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #606266;
  font-size: 15px;
  border-bottom: 1px solid #f5f5f5;
}
.pricing-features li:last-child { border-bottom: none; }

/* Footer */
.footer {
  background: #1a1a2e;
  color: #c0c4cc;
  padding: 60px 0 0;
}
.footer-section { margin-bottom: 32px; }
.footer-logo { color: #26A65B; font-size: 22px; font-weight: 800; letter-spacing: 2px; margin: 0 0 12px; }
.footer-desc { font-size: 14px; line-height: 1.7; color: #909399; }
.footer-section h4 { color: #fff; font-size: 15px; margin: 0 0 16px; font-weight: 600; }
.footer-links {
  list-style: none;
  padding: 0;
  margin: 0;
}
.footer-links li { margin-bottom: 10px; }
.footer-links a { color: #909399; text-decoration: none; font-size: 14px; transition: color 0.2s; }
.footer-links a:hover { color: #26A65B; }
.footer-contacts {
  list-style: none;
  padding: 0;
  margin: 0;
}
.footer-contacts li {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  color: #909399;
  font-size: 14px;
}
.footer-contacts .el-icon { color: #26A65B; font-size: 16px; }
.footer-bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  padding: 20px 0;
  margin-top: 20px;
  text-align: center;
}
.footer-bottom p { color: #606266; font-size: 13px; margin: 0; }

/* Register type modal */
.register-type-options {
  display: flex;
  gap: 20px;
}
.register-type-card {
  flex: 1;
  text-align: center;
  padding: 32px 20px;
  border: 2px solid #ebeef5;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.2s;
}
.register-type-card:hover {
  border-color: #26A65B;
  box-shadow: 0 4px 16px rgba(38, 166, 91, 0.15);
  transform: translateY(-2px);
}
.register-type-card h3 {
  margin: 16px 0 8px;
  color: #1a1a2e;
  font-size: 18px;
}
.register-type-card p {
  color: #909399;
  font-size: 13px;
  line-height: 1.5;
  margin: 0;
}
</style>
