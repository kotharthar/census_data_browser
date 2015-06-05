-- BaselineData_Census_2014_Provisional_Results_MIMU_10Oct2014.xlsx 
-- Provisional Results of the 2014 Myanmar Population and Housing Census													
-- Distribution of State and Region Population by Sex, Type of Place of Residence (Urban/Rural)													
-- State-Region Pop
CREATE TABLE state_pop(
    state_pcode text,
    state_name text,
    tp_both_sexes integer,
    tp_males integer,
    tp_females integer,
    ub_both_sexes integer,
    ub_males integer,
    ub_females integer,
    ru_both_sexes integer,
    ru_males integer,
    ru_females integer,
    pop_ratio real,
    mf_ratio real,
    remark text
);

-- BaselineData_Census_2014_Provisional_Results_MIMU_10Oct2014.xlsx 
-- Provisional Results of the 2014 Myanmar Population and Housing Census															
-- Distribution of Township and Sub-Township Enumerated Population by Sex, Type of Household and Administrative area (State/Region, District and Township	
-- TS+SubTS+Pop
CREATE TABLE ts_subts_pop(
    state_name text,
    district_name text,
    township_pcode text,
    township_name text,
    tp_both_sexes integer,
    tp_males integer,
    tp_females integer,
    cvh_households integer,
    cvh_both_sexes integer,
    cvh_males integer,
    cvh_females integer,
    cvh_hh_size real,
    ins_both_sexes integer,
    ins_males integer,
    ins_females integer,
    remark text
);

-- BaselineData_Census_2014_Provisional_Results_MIMU_10Oct2014.xlsx 
-- Provisional Results of the 2014 Myanmar Population and Housing Census									
-- Size of Enumerated Population and Households in Cities and State/ Region Capitals 									
-- State-Region City Pop
CREATE TABLE state_city_pop(
    state_pcode text,
    state_name text,
    city_pcode text,
    city_name text,
    household integer,
    total integer,
    male integer,
    female integer,
    hhsize real,
    mf_ratio real
);

-- IMPORT Command
.separator "\t"
.import state_pop.csv state_pop
.import ts_subts_pop.csv ts_subts_pop
.import state_city_pop.csv state_city_pop
