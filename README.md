# nCoV
获取2019-nCoV新型冠状病毒的实时信息
### 可以
- 指定省份的疫情信息
- 所有区域的疫情信息
- 疫情简报
- 搭建简单的API服务器
- 创建DUMP文件


具体使用方法请前往test文件夹查看具体使用方法

### API
返回JSON
1. **地址:** http://cxz.moe:1314/2019-nCoV/
```json
返回为空
```
2. **获取简单状态:** http://cxz.moe:1314/2019-nCoV/status
```json
{"confirmCount":2070,"suspectCount":2684,"deadCount":56,"hintWords":"中央应对疫情工作领导小组：将适当延长春节假期","recentTime":"2020-01-26 20:11","cure":49,"useTotal":true}
```
3. **获取所有区域的状态:** http://cxz.moe:1314/2019-nCoV/areas
```json
{"count":225,"data":[{"area":"湖北武汉","city":null,"confirm":41,"country":"中国","day":"1.12","dead":1,"district":null,"heal":null,"suspect":0,"time":""},{"area":"湖北武汉","city":null,"confirm":41,"country":"中国","day":"1.13","dead":1,"district":null,"heal":null,"suspect":0,"time":""},{"area":"湖北武汉","city":null,"confirm":41,"country":"中国","day":"1.14","dead":1,"district":null,"heal":null,"suspect":0,"time":""},{"area":"湖北武汉","city":null,"confirm":41,"country":"中国","day":"1.15","dead":2,"district":null,"heal":null,"suspect":0,"time":""},{"area":"湖北武汉","city":null,"confirm":45,"country":"中国","day":"1.16","dead":2,"district":null,"heal":null,"suspect":0,"time":""}]}
...

```
4. **获取指定区域的状态:** http://cxz.moe:1314/2019-nCoV/areas?area=湖北武汉
```json
{"area":"湖北武汉","confirmCount":270,"suspectCount":0,"deadCount":4,"heal":0}
```