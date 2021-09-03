clear all;

%location = '..\Evaluation\SpeechDetection_LOCAL_0001_0400\';
%file_path = fullfile(location, 'energyLog.csv');
%energy_data = readtable(file_path);

data = readtable('opp-edge-eval-runs.csv');

timestamp_datetime = datetime(data.time_stamp,'InputFormat','uuuu-MM-dd''T''HH:mm:ss.SSSSSS''+02:00''');
data.timestamp_datetime = timestamp_datetime;

providers = unique(data.provider_id);

for inx = 1:length(providers) 
    provider = providers{inx}
    
    data(strcmp(data.command,'Execution')==1 & strcmp(data.provider_id,provider)==1,:)
end

% executions = data(strcmp(data.command,'Execution')==1,:);
execution_start = executions(executions.step == 1,:);
execution_end = executions(executions.step == 2,:);

upper_bound = max(data.run_number);

results = ones(upper_bound+1,2);
results(:,1) = 0:upper_bound;

for i = 0:upper_bound

    start_time = execution_start.timestamp_datetime(execution_start.run_number == i);
    end_time = execution_end.timestamp_datetime(execution_end.run_number == i);
    duration = end_time - start_time;
    
    results(i+1,2) = milliseconds(duration);
    
    %disp(['Run number ' num2str(i) ': ' char(duration)]);  
end

figure()
plot(results(:,2));
mean_duration = mean(results(:,2));
hold on;
yline(mean_duration);





