<div style="text-align:center;"> 
<h1>buding-job(进行中)</h1>
</div>

<div style="text-align:center;">
<h2>简介</h2>
<h3>项目背景和说明</h3>
<p>本项目是基于Golang的轻量级分布式定时任务调度平台,同时也是基于作者原来的xll-job重构</p>
<p>重构的原因是原来的xll-job设计不合理，包的依赖混乱，其次是基于Robfig很久没有维护，同时也不适应当前设计</p>
<p>重构的项目作者除参考xxl-job之外还参考了PowerJob的设计理念,再加上自己的设计，取长补短</p>
<p>改进后的buding-job采用了更为底层的cronexpr库来解析cron表达式,同时自己实现了定时扫描的逻辑</p>
<p>项目在原来的基础上进行了很多的改动，如改进型了路由模式，重新梳理启动流程，改进Grpc调度，同时也优化了大量的代码</p>
<p>之所以重构的项目命名为buding-job,是为了纪念本人养的第一条宠物狗,名字叫布丁</p>
<img alt="img.png" height="100" src="static/img/img.png" width="100"/>
<h3>主要特性</h3>
<p> 操作简单(类似于xxl-job)，通过web界面实现定时任务配置</p>
<p> 支持单机、广播执行模式</p>
<p>支持CRON表达式、固定频率策略</p>
<p> 支持java/handler模式的执行，也可以通过sdk中grpc自行设计</p>
<p>支持多种数据库，依赖于Gorm</p>
<p> 高性能/低内存 基于编译型语言打造，直接运行二进制文件，无需依赖于jvm</p>
<p>提供重试/超时/故障转移等多种功能</p>

</div>

#### 同类产品对比

|        | QuartZ                         | xxl-job                          | BuDing-job                          |
|--------|--------------------------------|----------------------------------|-------------------------------------|
| 定时类型   | CRON                           | CRON                             | **CRON、固定频率**                       |
| 任务类型   | 内置Java                         | 内置Java、GLUE Java、Shell、Python等脚本 | **Java/JobHandle/Shell/基于Grpc自己定制** |
| 分布式计算  | 无                              | 静态分片                             | **(设计中)**                           |
| 集群部署   | 不支持                            | 每个节点都会扫描所有任务                     | **raft协议/智能分片/任务由分片单独管理**           |
| 在线任务治理 | 不支持                            | **支持**                           | **支持**                              |
| 日志白屏化  | 不支持                            | **支持**                           | **设计中**                             |
| 报警监控   | 无                              | **邮件**                           | **邮件**                              |
| 系统依赖   | JDBC支持的关系型数据库（MySQL、Oracle...） | mysql                            | **基于gorm,支持多种数据库**                  |
| DAG工作流 | 不支持                            | 不支持                              | **设计中**                             |
| 分布式调度  | 不支持                            | **支持**                           | **支持**                              |
| 运行占用内存 | 根据java服务决定                     | 500-n/mb(根据实时运行情况)               | **10-50/mb**                        |
| 性能     | 受限于jvm                         | 受限于jvm                           | **操作系统可执行文件,性能强劲**                  |
| 文件大小   | 根据java服务决定                     | 50mb左右                           | **10mb左右**                          |
| 运行依赖   | jvm                            | jvm                              | **二进制文件不依赖运行环境**                    |
| 在线监控   | 不支持                            | 不支持                              | **实时监控程序运行情况**                      |

#### 软件架构

软件架构说明

#### 安装教程

1. xxxx
2. xxxx
3. xxxx

#### 使用说明

1. xxxx
2. xxxx
3. xxxx

#### 参与贡献

1. Fork 本仓库
2. 新建 Feat_xxx 分支
3. 提交代码
4. 新建 Pull Request

