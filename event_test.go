package techlog

import (
	"reflect"
	"testing"
)

func Test_parseTechData(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []Event
	}{
		{"simple parser",
			args{data: []byte(t1)},
			[]Event{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTechData(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTechData() = %v, want %v", got, tt.want)
			}
		})
	}
}

var t2 = `01:33.503039-0,SDBL,6,process=rphost,p:processName=MyTempDbHost,OSThread=2424,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=1,Func=BeginTransaction,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ВыполнениеЗапросаSDBLНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 89 : НачатьТранзакцию();'`

var t1 = `29:40.227023-0,EXCP,4,process=rphost,p:processName=MyTempDbHost,OSThread=2840,t:clientID=12,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=4,SessionID=15,Usr=DefUser,AppID=1CV8C,Exception=580392e6-ba49-4280-ac67-fcd6f2180121,Descr='src\VResourceInfoBaseImpl.cpp(1129):
580392e6-ba49-4280-ac67-fcd6f2180121: Неспецифицированная ошибка работы с ресурсом
Ошибка при выполнении запроса POST к ресурсу /e1cib/logForm:
8d366056-4d5a-4d88-a207-0ae535b7d28e: Ошибка при вызове метода контекста (Выполнить)
{ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма(52)}:	Результат = Запрос.Выполнить();
f08d92f8-9eb2-4e19-9dd9-977d907cec2d
ae209c88-6b01-464c-adc9-0b72e240492f: {(4, 2)}: Таблица не найдена "Справочник.ИмяНесуществующегоСправочника"
<<?>>Справочник.ИмяНесуществующегоСправочника КАК КакойТоСправочник'
31:42.864025-1,SCRIPTCIRCREFS,5,process=rphost,

31:42.864025-1,SCRIPTCIRCREFS,3,process=rphost,p:processName=MyTempDbHost,OSThread=14024,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,ModuleName=Обработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма,ProcedureName=ЦиклическиеСсылкиНаСервере,Cycles='VariableName:
Структура1
CircularRefsMembers:
Структура1, Структура1.СсылкаНаСтруктуру2, Структура1
VariableName:
Структура2
CircularRefsMembers:
Структура2, Структура2.СсылкаНаСтруктуру1, Структура2
',Context='Форма.Вызов : Обработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ЦиклическиеСсылкиНаСервере
Обработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 79 : КонецПроцедуры'
01:33.503020-1,SDBL,4,process=rphost,p:processName=MyTempDbHost,OSThread=2424,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=0,Sdbl=SCHEMA SET GLOBAL,Rows=0
01:33.503022-1,SDBL,4,process=rphost,p:processName=MyTempDbHost,OSThread=2424,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=0,Sdbl=schema get local,Rows=0
01:33.503039-0,SDBL,6,process=rphost,p:processName=MyTempDbHost,OSThread=2424,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=1,Func=BeginTransaction,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ВыполнениеЗапросаSDBLНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 89 : НачатьТранзакцию();'

01:33.518001-14953,SDBL,6,process=rphost,p:processName=MyTempDbHost,OSThread=2424,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=1,Sdbl='SELECT
Q_000_T_001.ID
FROM
Reference47 Q_000_T_001
WHERE
(Q_000_T_001.Description  =  "Тест")

',Rows=1,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ВыполнениеЗапросаSDBLНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 101 : РезультатЗапроса = Запрос.Выполнить();'
49:44.366051-1,DBMSSQL,6,process=rphost,p:processName=MyTempDbHost,OSThread=24480,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=1,dbpid=54,Sql='SELECT
MAX(T1._Code)
FROM dbo._Reference47 T1',Rows=1,RowsAffected=-1,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ВыполнениеЗапросаSDBLНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 107 : НовыйОбъект.Записать();'

49:44.366068-1,DBMSSQL,6,process=rphost,p:processName=MyTempDbHost,OSThread=24480,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=1,dbpid=54,Sql='INSERT INTO dbo._Reference47 (_IDRRef,_Marked,_PredefinedID,_Code,_Description) VALUES(?,?,?,?,?)',Prm="
p_0: 0xB681D8F2CA20664B11EAE1169D5AADC6
p_1: FALSE
p_2: 0x00000000000000000000000000000000
p_3: '000000001'
p_4: 'Тест'
",RowsAffected=1,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ВыполнениеЗапросаSDBLНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 107 : НовыйОбъект.Записать();'

49:44.366070-1,DBMSSQL,6,process=rphost,p:processName=MyTempDbHost,OSThread=24480,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=1,dbpid=54,Sql='SELECT
T2._IDRRef,
T2._Version
FROM dbo._Reference47 T2
WHERE T2._IDRRef IN (?)
p_0: 0xB681D8F2CA20664B11EAE1169D5AADC6
',Rows=1,RowsAffected=-1,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ВыполнениеЗапросаSDBLНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 107 : НовыйОбъект.Записать();'

49:44.366074-1,DBMSSQL,6,process=rphost,p:processName=MyTempDbHost,OSThread=24480,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=1,dbpid=54,Sql="SELECT
T1._IDRRef
FROM dbo._Reference47 T1
WHERE T1._Code = ?
p_0: '000000001'
",Rows=1,RowsAffected=-1,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ВыполнениеЗапросаSDBLНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 107 : НовыйОбъект.Записать();'

49:44.366079-1,DBMSSQL,5,process=rphost,p:processName=MyTempDbHost,OSThread=24480,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=1,dbpid=54,Sql="{call sp_executesql(N'SELECT Creation,Modified,Attributes,DataSize,BinaryData FROM Params WHERE FileName = @P1 ORDER BY PartNo', N'@P1 nvarchar(128)', N'ibparams.inf')}",Rows=1,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ВыполнениеЗапросаSDBLНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 107 : НовыйОбъект.Записать();'

58:37.507025-1,DBMSSQL,6,process=rphost,p:processName=MyTempDbHost,OSThread=22564,t:clientID=45,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=11,SessionID=20,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Trans=1,dbpid=54,Sql="SELECT
T1._IDRRef
FROM dbo._Reference47 T1
WHERE (T1._Description = ?)
p_0: 'Тест'
",Rows=1,RowsAffected=-1,planSQLText='
1, 1, 1, 0.00313, 0.000158, 23, 0.00328, 1,   |--Index Seek(OBJECT:([MyTempDbHost].[dbo].[_Reference47].[_Reference47_3] AS [T1]), SEEK:([T1].[_Description]=[@P1]) ORDERED FORWARD)
',Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ВыполнениеЗапросаSDBLНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 101 : РезультатЗапроса = Запрос.Выполнить();'
11:15.624037-3,TLOCK,4,process=rphost,p:processName=MyTempDbHost,OSThread=19392,t:clientID=28,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=8,SessionID=23,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Regions=InfoRg49.DIMS,Locks='InfoRg49.DIMS Exclusive Fld50=47:b681d8f2ca20664b11eae1169d5aadc6',WaitConnections=,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ПростаяУправляемаяБлокировкаНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 144 : Блокировка.Заблокировать();'
11:15.624043-3,TLOCK,4,process=rphost,p:processName=MyTempDbHost,OSThread=19392,t:clientID=28,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=8,SessionID=23,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Regions=InfoRg49.DIMS,Locks='InfoRg49.DIMS Shared Fld50=47:b681d8f2ca20664b11eae1169d5aadc6',WaitConnections=,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ПростаяУправляемаяБлокировкаНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 148 : Набор.Прочитать();'
11:15.624051-3,TLOCK,4,process=rphost,p:processName=MyTempDbHost,OSThread=19392,t:clientID=28,t:applicationName=1CV8C,t:computerName=YY-COMP,t:connectID=8,SessionID=23,Usr=DefUser,AppID=1CV8C,DBMS=DBMSSQL,DataBase=YY-COMP\MyTempDbHost,Regions=InfoRg49.DIMS,Locks='InfoRg49.DIMS Exclusive Fld50=47:b681d8f2ca20664b11eae1169d5aadc6',WaitConnections=,Context='Форма.Вызов : ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Модуль.ПростаяУправляемаяБлокировкаНаСервере
ВнешняяОбработка.ТестированиеТехнологическогоЖурнала.Форма.Форма.Форма : 150 : Набор.Записать();'
`
