# buding-job

#### 介绍
为作者xll-job重构版,之所以重构是因为xll-job包之间的依赖很混乱<br/>
其次是因为使用了第三方robfig的cron库,这个库最近维护时间是三年前,不适合现在的逻辑，同时也很不方便<br/>
改进后的buding-job采用了更为底层的cronexpr库来解析cron表达式,同时自己实现了定时扫描的逻辑<br/>
之所以叫buding-job,是为了纪念本人养的第一条宠物狗,名字叫布丁
<br/>
<img alt="img.png" height="100" src="img/img.png" width="100"/>

### 同类产品对比
|                | QuartZ                   | xxl-job                          | BuDing-job              |
| -------------- | ------------------------ |----------------------------------|-------------------------|
| 定时类型       | CRON                     | CRON                             | CRON、固定频率               |
| 任务类型       | 内置Java                 | 内置Java、GLUE Java、Shell、Python等脚本 | **内置Java、可以根据SDk自己定制**  |
| 分布式计算     | 无                       | 静态分片                             | **MapReduce动态分片**       |
| 在线任务治理   | 不支持                   | 支持                               | **支持**                  |
| 日志白屏化     | 不支持                   | 支持                               | **支持**                  |
| 调度方式及性能 | 基于数据库锁，有性能瓶颈 | 基于数据库锁，有性能瓶颈                     | **无锁化设计，性能强劲无上限**       |
| 报警监控       | 无                       | 邮件                               | **WebHook、邮件、钉钉与自定义扩展** |
| 系统依赖       | JDBC支持的关系型数据库（MySQL、Oracle...）                    | **任意Spring Data Jpa支持的关系型数据库（MySQL、Oracle...）** |
| DAG工作流      | 不支持                   | 不支持                              | **支持**                  |

#### 软件架构
软件架构说明


#### 安装教程

1.  xxxx
2.  xxxx
3.  xxxx

#### 使用说明

1.  xxxx
2.  xxxx
3.  xxxx

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request

